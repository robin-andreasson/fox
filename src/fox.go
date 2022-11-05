package fox

//Later import net, errors
import (
	"fmt"
	"log"
	"net"

	parser "github.com/Robster0/http-parser"
)

type connection struct {
	Scheme string
}

type router struct {
	handlers []handler
	Port     int
}

func NewRouter() *router {
	return &router{}
}

func (r *router) Listen(port int, cb func(err error)) {
	ln, err := net.Listen("tcp", fmt.Sprint(":", port))

	r.Port = port

	if err != nil {
		cb(err)

		return
	}

	cb(nil)

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Panic(err)
		}

		go request(conn, *r)
	}
}

func request(conn net.Conn, r router) {

	var body []byte
	var c Context

	for {
		buffer := make([]byte, 1024*12)

		n, err := conn.Read(buffer)

		if err != nil {
			return
		}

		if len(c.Headers) == 0 {
			headers_bytes, body_bytes := parser.FirstInstance(buffer[0:n], "\r\n\r\n")

			c.Method, c.Url, c.Headers = parser.Headers(string(headers_bytes))

			if len(body_bytes) > 0 {
				body = append(body, body_bytes...)
			}

		} else {
			body = append(body, buffer[0:n]...)
		}

		if c.Headers["Content-Length"] == fmt.Sprint(len(body)) || c.Headers["Content-Length"] == "" {
			c.conn = conn
			r.handleRequests(c, body)
			break
		}
	}
}

func (r *router) handleRequests(c Context, raw []byte) {

	/*fmt.Println("HANDLEREQUESTS HAS BEEN CALLED")

	fmt.Println("Headers: ", c.Headers)
	fmt.Println("Raw Body: ", string(raw))

	fmt.Println("r: ", r)*/

	for i := 0; i < len(r.handlers); i++ {
		r.handlers[i].handler(c)
	}
}

var Test int = 2
