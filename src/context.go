package fox

import (
	"fmt"
	"net"
	"time"
)

type Context struct {
	Url        string
	Method     string
	Headers    map[string]string
	setHeaders map[string]string

	Body string
	Raw  []byte

	conn net.Conn
}

func (c *Context) Serve(nxtC Context) {

	//Epic things happening right here
}

func (c *Context) SetHeader(key string, value string) {
	c.setHeaders[key] = value
}

func json() {

	/*m := map[string]interface{}

	  err := json.Unmarshal([]byte(input), &m)
	  if err != nil {
	      panic(err)
	  }
	  fmt.Println(m)*/
}

func (c *Context) S(code int, body string) {

	c.conn.Write([]byte(c.setResponse(code, body)))

	c.conn.Close()
	//Send data
}

func (c *Context) setResponse(code int, body string) string {

	response := fmt.Sprint("HTTP/1.1 ", code, "\r\n")

	response += fmt.Sprint("Date: ", time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT"), "\r\n")

	response += "X-Powered-By: fox\r\n"

	response += "Connection: keep-alive\r\n"

	for key, value := range c.setHeaders {
		response += fmt.Sprint(key, ": ", value, "\r\n")
	}

	if body != "" {
		response += fmt.Sprint("Content-Length: ", len([]byte(body)), "\r\n")

		if c.Method != "HEAD" {
			response += fmt.Sprint("\r\n", body)
		}
	}

	return response
}
