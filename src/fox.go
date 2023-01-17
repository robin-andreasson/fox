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
	"strings"

	"github.com/robin-andreasson/fox/parser"
)

type router struct {
	handlers []handler
	prefix   string

	static map[string]static_content
}

type static_content struct {
	path string
	rex  string
}

func NewRouter() *router {
	return &router{}
}

/*
Get the value from nested map interfaces

type assertion shorthand

error occurs if a nested target is
*/
func Get(target any, keys ...string) any {

	if len(keys) == 0 {
		return target
	}

	key := keys[0]
	keys = keys[1:]
	t := reflect.TypeOf(target)

	if t == nil {
		log.Panic(errors.New("Can't nest key \"" + key + "\" because previous key was nil inside the target map"))
	}

	if t.Kind() != reflect.Map {
		log.Panic(errors.New(fmt.Sprint("target is not type of map but is ", reflect.TypeOf(target).Kind())))
	}

	next := target.(map[string]any)

	return Get(next[key], keys...)
}

/*
Statically serve files

"name" parameter is the name of the target folder

"relative_path" parameter is the path relative to the target folder, will use name if not specified.

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
		path = filepath.Join(call_path, "/"+name)
	} else {
		path = filepath.Join(call_path, relative_path[0])
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Panic("Cannot find the specified directory")
	}

	if r.static == nil {
		r.static = map[string]static_content{}
	}

	rex := `\/` + name + `\/.+`

	r.static[name] = static_content{filepath.Dir(path), rex}
}

func (r *router) Listen(port int) error {
	ln, err := net.Listen("tcp", fmt.Sprint(":", port))

	if err != nil {
		return err
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
	var temp_buffer []byte
	var c Context

	for {
		buffer := make([]byte, 1024)

		n, _ := conn.Read(buffer)

		temp_buffer = append(temp_buffer, buffer[0:n]...)

		if len(c.Headers) == 0 {

			headers_bytes, body_bytes := parser.FirstInstance(temp_buffer, "\r\n\r\n")

			if headers_bytes == nil {
				continue
			}

			//set context values
			c.Method, c.Url, c.Headers = parser.Headers(string(headers_bytes))
			c.setHeaders = make(map[string][]string)

			fmt.Println(c.Headers["Content-Type"])

			if len(body_bytes) > 0 {
				body = append(body, body_bytes...)
			}

		} else {
			body = append(body, buffer[0:n]...)
		}

		if c.Headers["Content-Length"] == fmt.Sprint(len(body)) || c.Method == "GET" {

			c._conn = conn
			r.handleRequests(c, body)

			break
		}
	}
}

func (r *router) handleRequests(c Context, body []byte) {

	//Checks if the url path is related to any of the static handlers
	if r.handleStatic(&c) {
		return
	}

	for _, handler := range r.handlers {

		if handler.method != c.Method {
			continue
		}

		match, queries, params := parser.Url(c.Url, handler.path, handler.rex, handler.params)

		if !match {
			continue
		}

		c.Params = params
		c.Query = queries

		handleBody(body, &c)

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
}

func handleBody(body []byte, c *Context) {
	if c.Headers["Content-Type"] == "" {
		return
	}

	segments := strings.Split(c.Headers["Content-Type"], "; ")

	switch segments[0] {
	case "application/json":
		c.Json = parser.JSON(body)
	case "multipart/form-data":

		delimiter := strings.Split(segments[1], "boundary=")[1]

		c.FormData = parser.FormData(body, []byte("--"+delimiter))
	case "application/x-www-form-urlencoded":

		c.Form = parser.Urlencoded(string(body))

	}
}

func (r *router) handleStatic(c *Context) bool {

	for _, s := range r.static {

		rex, err := regexp.Compile(s.rex)

		if err != nil {
			continue
		}

		if rex.FindString(c.Url) != "" {

			path := s.path + c.Url

			//If file does not exist return true anyways since it still matched the prefix
			if _, err := os.Stat(path); os.IsNotExist(err) {
				return true
			}

			//Send the file to the request endpoint
			c.File(Status.Ok, path)

			return true
		}
	}

	return false
}
