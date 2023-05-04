package fox

//Later import net, errors
import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/robin-andreasson/fox/parser"
)

type router struct {
	handlerstmpl
	static *map[string]static
}

type groupedrouter struct {
	handlerstmpl
}

type static struct {
	path string
	rex  string
}

/*
Initializes root router
*/
func Init() *router {
	return &router{handlerstmpl: handlerstmpl{handlers: &[]handler{}, prefix: ""}, static: &map[string]static{}}
}

/*
Create a group by specifying a path prefix
*/
func (r *router) Group(group string) *groupedrouter {

	if group == "" {
		log.Panic("unnecessary grouping")
	}

	if group[0] != '/' {
		group = "/" + group
	}

	nr := r.handlerstmpl

	nr.prefix += group

	return &groupedrouter{handlerstmpl: nr}
}

/*
Get value from nested map interfaces

error returns the zero value equivalent to type T
*/
func Get[T any](target any, keys ...string) T {
	targetType := reflect.TypeOf(target)

	if len(keys) == 0 {

		genericType := reflect.TypeOf(*new(T))

		if genericType != nil && targetType != genericType {
			return *new(T)
		}

		return target.(T)
	}

	key := keys[0]
	keys = keys[1:]

	if target == nil || targetType.Kind() != reflect.Map {
		return *new(T)
	}

	next := reflect.ValueOf(target).MapIndex(reflect.ValueOf(key))

	if next == reflect.Value(reflect.ValueOf(nil)) {
		return *new(T)
	}

	return Get[T](next.Interface(), keys...)
}

/*
Statically serve files

name is the name of the target directory

relative_path is the path relative to the target folder, will use name if not specified.

parameter is variadic but only allows one input as the purpose is only to make it optional
*/
func (r *router) Static(name string, relative_path ...string) {

	if len(relative_path) > 1 {
		log.Panic("only one relative_path argument is allowed")
	}

	_, call_path, _, _ := runtime.Caller(1)

	call_path = filepath.Dir(call_path)

	var path string

	if len(relative_path) == 0 {
		path = filepath.Join(call_path, "/"+name)
	} else {
		if relative_path[0][len(relative_path[0])-1] != '/' {
			relative_path[0] += "/"
		}

		path = filepath.Join(call_path, relative_path[0]+name)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Panic("could not find target directory " + name)
	}

	rex := `\/` + name + `\/.+`

	(*r.static)[name] = static{filepath.Dir(path), rex}
}

/*
Starts a server at the specified port
*/
func (r *router) Listen(port int) error {

	ln, err := net.Listen("tcp", fmt.Sprint(":", port))

	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			return err
		}

		go request(conn, *r)
	}
}

func request(conn net.Conn, r router) {

	var body []byte
	var header_buffer []byte
	var c Context

	var content_length int

	for {
		buffer := make([]byte, 4096)

		n, _ := conn.Read(buffer)

		if n == 0 {
			break
		}

		if len(c.Headers) == 0 {

			header_buffer = append(header_buffer, buffer[0:n]...)

			headers_bytes, body_bytes, found := parser.FirstInstance(header_buffer, "\r\n\r\n")

			if !found {
				continue
			}

			c.Method, c.Url, c.Headers = parser.Headers(string(headers_bytes))
			c.setHeaders = make(map[string][]string)
			c._conn = conn

			if c.Headers["Content-Length"] != "" {

				if cl, err := strconv.Atoi(c.Headers["Content-Length"]); err == nil {
					content_length = cl
				} else {
					c.Text(Status.BadRequest, "malformed content length")
					break
				}

			}

			if len(body_bytes) > 0 {
				body = append(body, body_bytes...)
			}

		} else {
			body = append(body, buffer[0:n]...)
		}

		if content_length == len(body) || c.Method == "GET" {
			r.handleRequests(c, body)
			break
		} else if content_length < len(body) {
			c.Text(Status.BadRequest, "malformed request syntax or not supported request technique/mechanism")
			break
		}
	}
}

func (r *router) handleRequests(c Context, body []byte) {

	for _, handler := range *r.handlerstmpl.handlers {
		if handler.method != c.Method {
			continue
		}

		match, queries, params := parser.Url(c.Url, handler.path, handler.rex, handler.params)

		if !match {
			continue
		}

		c.Raw = body
		c.Params = params
		c.Query = queries
		c.Cookies = parser.Cookies(c.Headers["Cookie"])

		handleBody(body, &c)
		handleSession(c.Cookies["FOXSESSID"], &c)
		handleRefresh(c.Headers["Authorization"], c.Cookies["FOXREFRESH"], &c)

		for _, function := range handler.stack {

			if err := function(&c); err != nil {
				c.Error = append(c.Error, err)
			}

			if !c._next {

				c._conn.Close()

				break
			}

			c._next = false
		}

		return
	}

	//Checks if the url path is related to any of the static handlers
	if !r.handleStatic(&c) {
		c.Status(Status.NotFound)
	}
}

func handleBody(body []byte, c *Context) {
	if c.Headers["Content-Type"] == "" {
		return
	}

	segments := strings.Split(c.Headers["Content-Type"], "; ")

	contentType := strings.ToLower(segments[0])

	switch contentType {

	case "application/json":
		if err := parser.JSONUnmarshal(string(body), &c.Body); err != nil {
			c.Error = append(c.Error, err)
			c.Body = make(map[string]any)
		}

	case "application/x-www-form-urlencoded":
		c.Body = parser.Urlencoded(string(body))
	case "multipart/form-data":
		if len(segments) <= 1 {
			c.Body = body
			return
		}

		delimiters := strings.Split(segments[1], "boundary=")

		if len(delimiters) <= 1 {
			c.Body = body
			return
		}

		delimiter := "--" + delimiters[1]

		c.Body = parser.FormData(body, []byte(delimiter)).(map[string]any)
	default:
		c.Body = body
	}
}

func (r *router) handleStatic(c *Context) bool {

	for _, s := range *r.static {

		rex, err := regexp.Compile(s.rex)

		if err != nil {
			continue
		}

		if rex.FindString(c.Url) != "" {

			path := s.path + c.Url

			if _, err := os.Stat(path); os.IsNotExist(err) {
				continue
			}

			err := c.File(Status.Ok, path)

			return err == nil
		}
	}

	return false
}
