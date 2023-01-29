package parser

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const rex_s = `[^\s"a-zA-Z0-9{},:[\]]+|"[^"\\]*(?:\\.[^"\\]*)*"|([a-zA-Z]+)|([0-9.]+)|(,|{|}|\[|\]|:)`

var keywords = map[string]map[string]bool{
	"{": {":": true},
	"[": {":": true},
	":": {},
	",": {"null": true, "false": true, "true": true, "}": true, "]": true},
}

var isOpening = map[string]bool{"{": true, "[": true}
var isClosing = map[string]bool{"}": true, "]": true}

func JSONUnmarshal(str string, output *any) error {

	rex := regexp.MustCompile(rex_s)

	segments := rex.FindAllString(str, -1)

	if len(segments) < 1 {
		return errors.New("nothing to parse")
	}

	opening := segments[0]

	if opening == "{" {
		*output = make(map[string]any)
	} else if opening == "[" {
		*output = []any{}
	} else {
		return errors.New("JSON has to start with an opening token (\"{\" or \"[\")")
	}

	if _, _, err := traverse(segments, &[]string{opening}, 1, output); err != nil {
		return err
	}

	return nil
}

func traverse(segments []string, stack *[]string, startIndex int, body *any) (int, any, error) {

	var mode string
	var name string

	if segments[startIndex-1] == "{" {
		mode = "object"
	} else if segments[startIndex-1] == "[" {
		mode = "array"
	}

	for i := startIndex; i < len(segments); i++ {

		seg := segments[i]
		previous := segments[i-1]
		keyword := keywords[seg]

		if isOpening[seg] {

			if name != "" {
				name = unicodedecode(name[1 : len(name)-1])
			}

			*stack = append(*stack, seg)

			next := nest(mode, seg, name, body)

			index, arr, err := traverse(segments, stack, i+1, next)

			if err != nil {
				return 0, nil, err
			}

			if arr != nil {
				setArrValue(mode, name, arr, body)
			}

			i = index

			continue

		} else if isClosing[seg] {

			latest := (*stack)[len(*stack)-1]

			if (latest == "{" && seg == "]") || (latest == "[" && seg == "}") {
				return 0, nil, errors.New("invalid closing scope, previous opening token was " + latest + " and current closing token is " + seg)
			}

			*stack = (*stack)[0 : len(*stack)-1]

			if IsArray(*body) {
				return i, *body, nil
			} else {
				return i, nil, nil
			}
		}

		if mode == "object" {

			if keyword == nil {

				if previous == ":" {

					value, valid := convertJsonValue(seg)

					if !valid {
						return 0, nil, errors.New("invalid value at " + seg)
					} else if name == "" {
						return 0, nil, errors.New("invalid name to value " + seg)
					}

					(*body).(map[string]any)[name[1:len(name)-1]] = value

					name = ""

				} else if previous == "{" || previous == "," {
					//name
					if !isStringLiteral(seg) {
						return 0, nil, errors.New("invalid name at " + seg)
					}

					name = seg
				} else {
					return 0, nil, errors.New("latest token is invalid")
				}

			} else if !keyword[previous] {

				if _, valid := convertJsonValue(previous); !valid {
					return 0, nil, errors.New("invalid name or value at " + seg)
				}
			}

			continue
		}

		if seg == "," {

			if _, valid := convertJsonValue(previous); !valid && previous != "}" && previous != "]" {
				return 0, nil, errors.New("invalid syntax inside array")
			}

			continue

		}

		if previous != "," && previous != "[" {
			return 0, nil, errors.New("invalid previous token")
		}

		value, valid := convertJsonValue(seg)

		if !valid {
			return 0, nil, errors.New("invalid value inside array")
		}

		*body = append((*body).([]any), value)
	}

	if len(*stack) != 0 {
		return 0, nil, errors.New("missing closing token for opening token " + (*stack)[len(*stack)-1])
	}

	return 0, nil, nil
}

func isStringLiteral(seg string) bool {
	first := seg[0]
	last := seg[len(seg)-1]

	if first != '"' || last != '"' {
		return false
	}

	return true
}

func convertJsonValue(seg string) (any, bool) {
	//If null then return nil since thats the correct equivalent
	if seg == "null" {
		return nil, true
	}

	if seg == "true" {
		return true, true
	}

	if seg == "false" {
		return false, true
	}

	if isStringLiteral(seg) {
		value := unicodedecode(seg[1 : len(seg)-1])

		value = strings.Replace(value, `"`, `\"`, -1)

		return value, true
	}

	if i, isNumber := getNumber(seg); isNumber {
		return i, true
	}

	return nil, false
}

func setArrValue(mode string, name string, arr any, body *any) {
	if mode == "array" {
		(*body).([]any)[len((*body).([]any))-1] = arr
	} else {
		(*body).(map[string]any)[name] = arr
	}
}

func nest(mode string, token string, name string, body *any) *any {

	var nest any
	var next any

	if token == "[" {
		nest = []any{}
	} else {
		nest = make(map[string]any)
	}

	if mode == "array" {
		*body = append((*body).([]any), nest)

		next = (*body).([]any)[len((*body).([]any))-1]

	} else {
		(*body).(map[string]any)[name] = nest

		next = (*body).(map[string]any)[name]
	}

	return &next
}

func JSONMarshal(v any) (string, error) {
	if !IsMap(v) && !IsArray(v) {
		value, err := convertGoValues(v)

		if err != nil {
			return "", err
		}

		return value, nil
	}

	s := ""

	if IsMap(v) {
		s += "{"

		comma := ""

		for name, value := range v.(map[string]any) {

			var next any

			if IsMap(value) {
				next = value.(map[string]any)
			} else {
				next = value
			}

			result, err := JSONMarshal(next)

			if err != nil {
				return "", err
			}

			s += comma +
				`"` + name + `"` +
				":" +
				result

			comma = ","
		}

		s += "}"
	} else {
		if isBytes(v) {
			b64str := b64.StdEncoding.EncodeToString(v.([]byte))

			b64str = "b64str"

			return b64str, nil
		}

		s += "["

		comma := ""

		for _, value := range v.([]any) {

			if IsMap(value) {
				result, err := JSONMarshal(value)

				if err != nil {
					return "", err
				}

				value = result
			} else {
				result, err := convertGoValues(value)

				if err != nil {
					return "", err
				}

				value = result
			}

			s += comma + value.(string)

			comma = ","
		}

		s += "]"
	}

	return s, nil
}

func convertGoValues(v any) (string, error) {

	if v == nil {
		return "null", nil
	}

	if reflect.TypeOf(v).Kind() == reflect.String {
		return `"` + v.(string) + `"`, nil
	}

	if v == true {
		return "true", nil
	}

	if v == false {
		return "false", nil
	}

	if reflect.TypeOf(v).Kind() == reflect.Int {
		return fmt.Sprintf("%v", v), nil
	}

	if reflect.TypeOf(v).Kind() == reflect.Float64 {
		return fmt.Sprintf("%v", v), nil
	}

	return "", errors.New("invalid json value")
}
