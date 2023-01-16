package fox

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"encoding/base64"
)

type Context struct {
	Url        string
	Method     string
	Headers    map[string]string
	setHeaders map[string][]string

	Json     interface{}
	Form     interface{}
	FormData map[string]interface{}

	Params map[string]string
	Query  map[string]string
	Raw    []byte

	_next bool
	_conn net.Conn
}

type CookieAttributes struct {
	HttpOnly  bool
	BASE64    bool
	Secure    bool
	Path      string
	Domain    string
	SameSite  string //strict, lax or none are the only accepted values
	ExpiresIn int
}

func (c *Context) Nested(target map[string]interface{}, keys ...string) interface{} {

	fmt.Println(target)
	if len(keys) == 0 {
		return target
	}

	key := keys[0]
	keys = keys[1:]

	return c.Nested(target[key].(map[string]interface{}), keys...)
}

func (c *Context) Next() {
	c._next = true
}

func (c *Context) SetHeader(key string, value string) {
	c.setHeaders[key] = append(c.setHeaders[key], value)
}

/*func json() {

	m := map[string]interface{}

	  err := json.Unmarshal([]byte(input), &m)
	  if err != nil {
	      panic(err)
	  }
	  fmt.Println(m)
}*/

func (c *Context) Head(code int) {
	if c.Method != "HEAD" {
		log.Panic("head http response function should only be called during a response when the request method is 'HEAD'")
	}

	if err := c.response(code, []byte{}); err != nil {
		log.Panic(err)
	}
}

func (c *Context) String(code int, body string) {
	if err := c.response(code, []byte(body)); err != nil {
		log.Panic(err)
	}
}

func (c *Context) File(code int, path string) {

	buffer, err := os.ReadFile(path)

	if err != nil {
		log.Panic(err)
	}

	if err = c.response(code, buffer); err != nil {
		log.Panic(err)
	}
}

func (c *Context) Status(code int) {
	if err := c.response(code, []byte{}); err != nil {
		log.Panic(err)
	}
}

/*
Set a cookie
*/
func (c *Context) Cookie(name string, value string, attributes CookieAttributes) {

	cookie := name + "="

	if attributes.BASE64 {
		cookie += base64.StdEncoding.EncodeToString([]byte(value))
	} else {
		cookie += value
	}

	if attributes.ExpiresIn != 0 {
		cookie += "; Expires=" + formatTime(time.Duration(attributes.ExpiresIn))
	}

	if attributes.HttpOnly {
		cookie += "; HttpOnly"
	}

	if attributes.Secure {
		cookie += "; Secure"
	}

	if attributes.Domain != "" {
		cookie += "; Domain=" + attributes.Domain
	}

	if attributes.Path != "" {
		cookie += "; Path=" + attributes.Path
	}

	switch strings.ToLower(attributes.SameSite) {
	case "strict":
		cookie += "; SameSite=Strict"
		break
	case "lax":
		cookie += "; SameSite=Lax"
		break
	case "none":
		cookie += "; SameSite=None"
		break
	case "":
		break
	default:
		log.Panic(fmt.Errorf("samesite attribute can only be the values 'strict', 'lax' and 'none'"))
	}

	c.SetHeader("Set-Cookie", cookie)
}

func (c *Context) response(code int, body []byte) error {

	response := fmt.Sprint("HTTP/1.1 ", code, "\r\n")

	response += fmt.Sprint("Date: ", formatTime(0), "\r\n")

	response += "X-Powered-By: fox\r\n"

	response += "Connection: keep-alive\r\n"

	for key, values := range c.setHeaders {
		for _, value := range values {
			response += fmt.Sprint(key, ": ", value, "\r\n")
		}
	}

	response_bytes := []byte(response)

	if len(body) > 0 {

		contentLength := []byte(fmt.Sprint("Content-Length: ", len(body), "\r\n\r\n"))

		response_bytes = append(response_bytes, contentLength...)

		response_bytes = append(response_bytes, body...)
	}

	_, err := c._conn.Write(response_bytes)

	c._conn.Close()

	if err != nil {
		return err
	}

	return nil
}
