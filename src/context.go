package fox

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

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
returns set response headers
*/
func (c *Context) ResHeaders() map[string][]string {
	return c.setHeaders
}

/*
Set a header by passing a name and value
*/
func (c *Context) SetHeader(name string, value string) error {

	name = strings.Title(strings.ToLower(name))

	if name == "Set-Cookie" {

		if len([]byte(value)) > 4093 {
			return errors.New("set-cookie value exceeded the size limit of 4093")
		}

		c.setHeaders[name] = append(c.setHeaders[name], value)
	} else {
		c.setHeaders[name] = []string{value}
	}

	return nil
}

func (c *Context) Head(code int) error {
	if c.Method != "HEAD" {
		return errors.New("head function can only be called when method is HEAD")
	}

	if err := c.response(code, []byte{}); err != nil {
		return err
	}

	return nil
}

/*
Send text back to the request endpoint

content type is set to text/html; charset=utf-8
*/
func (c *Context) Text(code int, body string) error {

	c.SetHeader("Content-Type", "text/html; charset=utf-8")

	return c.response(code, []byte(body))
}

/*
Send back file data to the request endpoint

basic mime types like images, zips, fonts, pdf and mp4 files are calculated.

mime types from script files that is in need for a sniffing technique is found through file extension
*/
func (c *Context) File(code int, path string) error {

	bytes, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	mime := parser.Mime(path, bytes)

	c.SetHeader("Content-Type", mime)

	return c.response(code, bytes)
}

/*
Send application/json response to the request endpoint
*/
func (c *Context) JSON(status int, body any) error {

	if !parser.IsMap(body) && !parser.IsArray(body) {
		return errors.New("invalid type for body, expected map or array/slice")
	}

	s, err := parser.JSONMarshal(body)

	if err != nil {
		return err
	}

	c.SetHeader("Content-Type", "application/json")

	return c.response(status, []byte(s))
}

func (c *Context) Status(status int) error {
	return c.response(status, []byte{})
}

func (c *Context) Redirect(path string) error {

	c.SetHeader("Location", path)

	return c.Status(Status.SeeOther)
}

/*
Set a cookie
*/
func (c *Context) Cookie(name string, value string, attributes CookieAttributes) error {

	cookie := name + "="

	if attributes.BASE64 {
		cookie += base64.StdEncoding.EncodeToString([]byte(value))
	} else {
		cookie += value
	}

	if attributes.MaxAge != 0 {
		cookie += "; Max-Age=" + fmt.Sprint(attributes.MaxAge)
	} else if attributes.ExpiresIn != 0 {
		cookie += "; Expires=" + formatTime(attributes.ExpiresIn)
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
		return errors.New("samesite attribute is limited between the values Strict, Lax and None")
	}

	return c.SetHeader("Set-Cookie", cookie)
}

func (c *Context) response(code int, body []byte) error {

	status := handleCors(c)

	if status != 0 {
		code = status
		body = []byte{}
	}

	response := fmt.Sprint("HTTP/1.1 ", code, "\r\n")

	response += fmt.Sprint("Date: ", formatTime(0), "\r\n")

	response += "X-Powered-By: fox\r\n"

	response += "Connection: keep-alive\r\n"

	for key, values := range c.setHeaders {
		for _, value := range values {
			response += key + ": " + value + "\r\n"
		}
	}

	response += "Content-Length: " + fmt.Sprint(len(body)) + "\r\n\r\n"

	response_bytes := []byte(response)

	if c.Method != "HEAD" {
		response_bytes = append(response_bytes, body...)
	}

	_, err := c._conn.Write(response_bytes)

	c._conn.Close()

	return err
}
