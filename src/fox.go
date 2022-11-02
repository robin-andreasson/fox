package fox

//Later import net, errors

type Context struct {
	Request  string
	Response string

	Serve func(c Context)
}

type handler struct {
	path    string
	method  string
	mw      []func(c Context)
	handler func(c Context)
}

type router struct {
	handlers []handler
}

func NewRouter() *router {
	var r router

	return &r
}

var Test Context = Context{"Request", "Response", nil}
