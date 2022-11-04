package fox

//Later import net, errors
import (
	"bufio"
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
}

func NewRouter() *router {
	return &router{}
}

func (r *router) Listen(port int) {
	ln, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Panic(err)
		}

		scanner := bufio.NewScanner(conn)

		fmt.Println(r)
		//var c Context

		var message []byte

		for {
			if !scanner.Scan() {
				r.handleRequests(message)

				break
			}

			message = append(message, scanner.Bytes()...)
		}

		defer conn.Close()

	}
}

func (r *router) handleRequests(data []byte) {

	headers_bytes, body_bytes := parser.FirstInstance(data, "\r\n\r\n")

	parser.Headers(string(headers_bytes))

	if len(body_bytes) > 0 {
		//Parse body
	}

}
