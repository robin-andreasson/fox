package fox

func (r *router) GET(path string, functions ...func(c Context)) {
	r.handlers = append(r.handlers, addHandler(path, "GET", functions))
}

func (r *router) POST(path string, functions ...func(c Context)) {
	r.handlers = append(r.handlers, addHandler(path, "POST", functions))
}

func (r *router) PUT(path string, functions ...func(c Context)) {
	r.handlers = append(r.handlers, addHandler(path, "PUT", functions))
}

func (r *router) DELETE(path string, functions ...func(c Context)) {
	r.handlers = append(r.handlers, addHandler(path, "DELETE", functions))
}

func (r *router) HEAD(path string, functions ...func(c Context)) {
	r.handlers = append(r.handlers, addHandler(path, "HEAD", functions))
}

func (r *router) PATCH(path string, functions ...func(c Context)) {
	r.handlers = append(r.handlers, addHandler(path, "PATCH", functions))
}

func addHandler(_path string, _method string, functions []func(c Context)) handler {
	lastIndex := len(functions) - 1

	var stack []func(c Context)

	for i := 0; i < lastIndex; i++ {
		stack = append(stack, functions[i])
	}

	return handler{path: _path, method: _method, mw: stack, handler: functions[lastIndex]}
}
