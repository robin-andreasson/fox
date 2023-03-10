package fox

import (
	"log"
	"regexp"
	"strings"
)

type handlerstmpl struct {
	handlers *[]handler
	prefix   string
}

type handler struct {
	path   string
	method string
	stack  []func(c *Context) error
	rex    string
	params [][]string
}

func (r *handlerstmpl) Get(path string, stack ...func(c *Context) error) {
	*r.handlers = append(*r.handlers, r.addHandler(path, "GET", stack))
}

func (r *handlerstmpl) Post(path string, stack ...func(c *Context) error) {
	*r.handlers = append(*r.handlers, r.addHandler(path, "POST", stack))
}

func (r *handlerstmpl) Put(path string, stack ...func(c *Context) error) {
	*r.handlers = append(*r.handlers, r.addHandler(path, "PUT", stack))
}

func (r *handlerstmpl) Delete(path string, stack ...func(c *Context) error) {
	*r.handlers = append(*r.handlers, r.addHandler(path, "DELETE", stack))
}

func (r *handlerstmpl) Head(path string, stack ...func(c *Context) error) {
	*r.handlers = append(*r.handlers, r.addHandler(path, "HEAD", stack))
}

func (r *handlerstmpl) Patch(path string, stack ...func(c *Context) error) {
	*r.handlers = append(*r.handlers, r.addHandler(path, "PATCH", stack))
}

func (r *handlerstmpl) Options(path string, stack ...func(c *Context) error) {
	*r.handlers = append(*r.handlers, r.addHandler(path, "OPTIONS", stack))
}

func (r *handlerstmpl) addHandler(path string, method string, stack []func(c *Context) error) handler {

	if len(stack) == 0 {
		log.Panic("one handler is required")
	}

	rex := regexp.MustCompile("^:([^;]+);(.+?)$|^:([^;]+)$")

	var paramArr [][]string
	params := map[string]bool{}

	path = r.prefix + path

	path_segs := strings.Split(path, "/")

	path_rex := regexp.QuoteMeta(path)

	for _, path_seg := range path_segs {

		found := rex.FindString(path_seg)

		if found == "" {

			emptypattern_rex := regexp.MustCompile("^:.+?;$")

			if emptypattern_rex.FindString(path_seg) != "" {
				log.Panic("regex pattern at " + path + " on param " + path_seg + " can not be empty")
			}

			temp := strings.ReplaceAll(regexp.QuoteMeta(path_seg), `\*`, ".*")
			path_rex = strings.ReplaceAll(path_rex, regexp.QuoteMeta(path_seg), temp)

			continue
		}

		//Cuts the first instance of ;
		param_name, param_pattern, hasRex := strings.Cut(path_seg[1:], ";")

		param_info := []string{param_name}

		if hasRex {
			param_info = append(param_info, "^"+param_pattern+"$")
		}

		paramArr = append(paramArr, [][]string{param_info}...)

		if params[param_name] {
			log.Panic("duplicate path params " + param_name + " at " + path)
		}

		params[param_name] = true

		path_rex = strings.Replace(path_rex, regexp.QuoteMeta(path_seg), "(.+?)", 1)
	}

	return handler{path: path, method: method, stack: stack, rex: "^" + path_rex + "$", params: paramArr}
}
