package fox

type handler struct {
	path    string
	method  string
	mw      []func(c Context)
	handler func(c Context)
}

func (r *router) GET(path string, functions ...func(c Context)) {
	r.addHandler(path, "GET", functions)
}

func (r *router) POST(path string, functions ...func(c Context)) {
	r.addHandler(path, "POST", functions)
}

func (r *router) PUT(path string, functions ...func(c Context)) {
	r.addHandler(path, "PUT", functions)
}

func (r *router) DELETE(path string, functions ...func(c Context)) {
	r.addHandler(path, "DELETE", functions)
}

func (r *router) HEAD(path string, functions ...func(c Context)) {
	r.addHandler(path, "HEAD", functions)
}

func (r *router) PATCH(path string, functions ...func(c Context)) {
	r.addHandler(path, "PATCH", functions)
}

func (r *router) addHandler(path string, method string, functions []func(c Context)) {
	lastIndex := len(functions) - 1

	var stack []func(c Context)

	for i := 0; i < lastIndex; i++ {
		stack = append(stack, functions[i])
	}

	r.handlers = append(r.handlers, handler{path: path, method: method, mw: stack, handler: functions[lastIndex]})
}
