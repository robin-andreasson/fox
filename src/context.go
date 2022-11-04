package fox

type Context struct {
	Request  string
	Response string

	Headers string

	statusCode int
	body       string
}

func (c *Context) Serve(nxtC Context) {

	//Epic things happening right here
}

func (c *Context) Status(code int) *Context {
	//Set status code

	c.statusCode = code

	return c
}

func (c *Context) Send(body string) {

	//Send data
	c.body = body
}
