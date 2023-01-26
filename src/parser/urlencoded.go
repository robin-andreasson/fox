package parser

import (
	"reflect"
	"regexp"
	"strings"
)

const (
	rex_s_str = "&|="
	rex_n_str = `\[(.+?)\]`
)

func Urlencoded(s string) map[string]any {
	body := make(map[string]any)

	rex_s := regexp.MustCompile(rex_s_str)
	rex_n := regexp.MustCompile(rex_n_str)

	seg := rex_s.Split(s, -1)

	for i := 0; i < len(seg); i += 2 {

		name := urldecode(seg[i])
		value := urldecode(seg[i+1])

		nestedkeys := rex_n.FindAllStringSubmatch(name, -1)

		//if there are nested keys
		if len(nestedkeys) != 0 {

			n, _, _ := strings.Cut(name, "[")

			body[n] = make(map[string]any)
			next := body[n]

			nested(nestedkeys, value, &next)
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
		(*body).(map[string]any)[name] = make(map[string]any)
	}

	next := (*body).(map[string]any)[name]

	nested(names, value, &next)
}
