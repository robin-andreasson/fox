package fox

//Later import net, errors
import (
	"errors"
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
	handlers *[]handler
	prefix   string

	preflight *handler

	static *map[string]static
}

type static struct {
	path string
	rex  string
}

/*
Initializes root router
*/
func Init() *router {
	return &router{handlers: &[]handler{}, preflight: &handler{}, static: &map[string]static{}}
}

/*
Create a group by specifying a path prefix
*/
func (r *router) Group(group string) *router {

	if group == "" {
		log.Panic("unnecessary grouping")
	}

	if group[0] != '/' {
		group = "/" + group
	}

	return &router{handlers: r.handlers, preflight: r.preflight, prefix: r.prefix + group, static: r.static}
}

/*
Get the value from nested map interfaces

error occurs if the target is nil or not a map
*/
func Get[T any](target any, keys ...string) T {

	if len(keys) == 0 {
		return target.(T)
	}

	key := keys[0]
	keys = keys[1:]

	t := reflect.TypeOf(target)

	if target == nil || t.Kind() != reflect.Map {
		log.Panic(errors.New("cannot nest target at key '" + key + "' because target is either not a map or nil"))
	}

	next := reflect.ValueOf(target).MapIndex(reflect.ValueOf(key)).Interface()

	return Get[T](next, keys...)
}

/*
Statically serve files

name is the name of the target directory

relative_path is the path relative to the target folder, will use name if not specified.

parameter is variadic but only allows one input as the purpose is only to make it optional
*/
func (r *router) Static(name string, relative_path ...string) {

	if len(relative_path) > 1 {
		log.Panic(errors.New("only one relative_path argument is allowed"))
	}

	_, call_path, _, _ := runtime.Caller(1)

	call_path = filepath.Dir(call_path)

	var path string

	if len(relative_path) == 0 {
		path = filepath.Join(call_path, r.prefix+"/"+name)
	} else {
		path = filepath.Join(call_path, r.prefix+relative_path[0])
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Panic("cannot find target directory " + name)
	}

	rex := `\/` + name + `\/.+`

	(*r.static)[name] = static{filepath.Dir(path), rex}
}

/*
Starts a server at the specified port
*/
func Listen(port int, r *router) error {

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
		buffer := make([]byte, 65535)

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

	for _, handler := range *r.handlers {
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

		for _, function := range handler.stack {

			function(&c)

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

	switch segments[0] {
	case "application/json":
		if err := parser.JSONUnmarshal(string(body), &c.Body); err != nil {
			c.Body = make(map[string]any)
		}

	case "application/x-www-form-urlencoded":
		c.Body = parser.Urlencoded(string(body))
	case "multipart/form-data":
		delimiter := strings.Split(segments[1], "boundary=")[1]

		c.Body = parser.FormData(body, []byte("--"+delimiter)).(map[string]any)
	default:
		c.Body = make(map[string]any)
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

			//Send the file to the request endpoint
			c.File(Status.Ok, path)

			return true
		}
	}

	return false
}
