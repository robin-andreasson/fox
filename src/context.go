package fox

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"encoding/base64"

	"github.com/robin-andreasson/fox/parser"
)

type Context struct {
	Url        string
	Method     string
	Headers    map[string]string
	setHeaders map[string][]string

	Body any

	Params  map[string]string
	Query   map[string]string
	Cookies map[string]string

	Raw []byte

	_next bool
	_conn net.Conn
}

type CookieAttributes struct {
	HttpOnly    bool
	BASE64      bool
	Secure      bool
	Partitioned bool
	Path        string
	Domain      string
	SameSite    string //strict, lax or none
	ExpiresIn   int
	MaxAge      int
}

func (c *Context) Next() {
	c._next = true
}

/*
Set a header by passing a name and value
*/
func (c *Context) SetHeader(name string, value string) {

	if strings.ToLower(name) == "set-cookie" {

		if len([]byte(value)) > 4093 {
			log.Panic("Set-Cookie value exceeded the size limit of 4093")
		}

		c.setHeaders[name] = append(c.setHeaders[name], value)
	} else {

		c.setHeaders[name] = []string{value}
	}
}

func (c *Context) Head(code int) {
	if c.Method != "HEAD" {
		log.Panic("'Head' function should only be called during http requests where the method is HEAD")
	}

	if err := c.response(code, []byte{}); err != nil {
		log.Panic(err)
	}
}

/*
Send text back to the request endpoint

content type is set to text/html; charset=utf-8
*/
func (c *Context) Text(code int, body string) {

	c.SetHeader("Content-Type", "text/html; charset=utf-8")

	if err := c.response(code, []byte(body)); err != nil {
		log.Panic(err)
	}
}

/*
Send back file data to the request endpoint

basic mime types like images, zips, fonts, pdf and mp4 files are calculated.

mime types from script files that is in need for a sniffing technique is found through file extension
*/
func (c *Context) File(code int, path string) {

	bytes, err := os.ReadFile(path)

	if err != nil {
		log.Panic(err)
	}

	mime := parser.Mime(path, bytes)

	c.SetHeader("Content-Type", mime)

	if err = c.response(code, bytes); err != nil {
		log.Panic(err)
	}
}

/*
Send application/json response to the request endpoint
*/
func (c *Context) JSON(status int, body any) {

	if !parser.IsMap(body) && !parser.IsArray(body) {
		log.Panic("invalid type for body, expected map or array/slice")
	}

	s, err := parser.JSONMarshal(body)

	if err != nil {
		log.Panic(err)
	}

	c.SetHeader("Content-Type", "application/json")

	if err := c.response(status, []byte(s)); err != nil {
		log.Panic(err)
	}
}

func (c *Context) Status(status int) {
	if err := c.response(status, []byte{}); err != nil {
		log.Panic(err)
	}
}

func (c *Context) Redirect(path string) {

	c.SetHeader("Location", path)

	c.Status(Status.SeeOther)
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

	if attributes.MaxAge != 0 {
		cookie += "; Max-Age=" + fmt.Sprint(attributes.MaxAge)
	} else if attributes.ExpiresIn != 0 {
		cookie += "; Expires=" + formatTime(time.Duration(attributes.ExpiresIn))
	}

	if attributes.Partitioned {
		cookie += "; Partitioned"
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
		log.Panic(fmt.Errorf("samesite attribute can only have the values 'Strict', 'Lax' and 'None'"))
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
			response += key + ": " + value + "\r\n"
		}
	}

	response_bytes := []byte(response)

	if c.Method != "HEAD" {

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
