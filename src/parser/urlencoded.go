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

			k, _, _ := strings.Cut(name, "[")

			if body[k] == nil || reflect.TypeOf(body[k]).Kind() != reflect.Map {
				body[k] = make(map[string]any)
			}

			body[k] = getNestedKeys(nestedkeys, value, body[k])

			continue
		}

		body[name] = value
	}

	return body
}

func getNestedKeys(names [][]string, value string, body any) any {
	if len(names) == 0 {
		return convertValue(value)
	}

	name := names[0][1]
	names = names[1:]

	next := body.(map[string]any)[name]

	if next == nil || reflect.TypeOf(next).Kind() != reflect.Map {
		next = make(map[string]any)
	}

	nested_value := getNestedKeys(names, value, next)

	body.(map[string]any)[name] = nested_value

	return body
}

func convertValue(v string) any {

	if i, isNumber := getNumber(v); isNumber {
		return i
	}

	return v
}
