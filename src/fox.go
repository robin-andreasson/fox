package fox

//Later import net, errors
import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/robin-andreasson/fox/parser"
)

type router struct {
	handlers []handler
	prefix   string
	Port     int
}

func NewRouter() *router {
	return &router{}
}

func (r *router) Listen(port int) error {
	ln, err := net.Listen("tcp", fmt.Sprint(":", port))

	r.Port = port

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

func (r *router) handleRequests(c Context, raw []byte) {

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

		if c.Headers["Content-Type"] != "" {
			switch c.Headers["Content-Type"] {
			case "application/json":
				err := json.Unmarshal(raw, &c.Json)

				if err != nil {
					log.Panic("error parsing json object")
				}
			}
		}

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
