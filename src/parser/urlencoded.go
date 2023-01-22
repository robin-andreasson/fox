package parser

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const (
	rex_segments = "&|="
	rex_nest     = `\[(.+?)\]`
)

func Urlencoded(s string) map[string]any {
	body := make(map[string]any)

	rex_s := regexp.MustCompile(rex_segments)
	rex_n := regexp.MustCompile(rex_nest)

	seg := rex_s.Split(s, -1)

	for i := 0; i < len(seg); i += 2 {

		name := urldecode(seg[i])
		value := urldecode(seg[i+1])

		nestedobj := rex_n.FindAllStringSubmatch(name, -1)

		if len(nestedobj) != 0 {

			n, _, _ := strings.Cut(name, "[")

			body[n] = make(map[string]any)
			next := body[n]

			nested(nestedobj, value, &next)
			continue
		}

		body[name] = value
	}

	return body
}

func nested(names [][]string, value string, body *any) {

	name := names[0][1]
	names = names[1:]

	t := (*body).(map[string]any)[name]

	if t == nil || reflect.TypeOf(t).Kind() != reflect.Map {
		fmt.Println(*body)
		(*body).(map[string]any)[name] = make(map[string]any)
	}

	next := (*body).(map[string]any)[name]

	if len(names) == 0 {
		(*body).(map[string]any)[name] = value
		return
	}

	nested(names, value, &next)
}
