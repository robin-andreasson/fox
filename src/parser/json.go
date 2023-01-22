package parser

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
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

func JSON(str string, output *any) error {

	rex := regexp.MustCompile(rex_s)

	segments := rex.FindAllString(str, -1)

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
				name = name[1 : len(name)-1]
			}

			*stack = append(*stack, seg)

			next := next_nestType(mode, seg, name, body)

			index, arr, err := traverse(segments, stack, i+1, next)

			if err != nil {
				return 0, nil, err
			}

			if arr != nil {
				setValue(mode, name, arr, body)
			}

			i = index

			continue

		} else if isClosing[seg] {

			latest := (*stack)[len(*stack)-1]

			if (latest == "{" && seg == "]") || (latest == "[" && seg == "}") {
				return 0, nil, errors.New("invalid closing scope, previous opening token was " + latest + " and current closing token is " + seg)
			}

			*stack = (*stack)[0 : len(*stack)-1]

			if reflect.TypeOf(*body).Kind() == reflect.Slice ||
				reflect.TypeOf(*body).Kind() == reflect.Array {
				return i, *body, nil
			} else {
				return i, nil, nil
			}

		}

		if mode == "object" {

			if keyword == nil {

				if previous == ":" {

					value, valid := Value(seg)

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
					return 0, nil, errors.New("invalid syntax at " + seg)
				}

			} else if !keyword[previous] {

				if _, valid := Value(previous); !valid {
					return 0, nil, errors.New("invalid syntax at " + seg + ", previous keyword is wrong (" + previous + ")")
				}
			}

			continue
		}

		if seg == "," {

			if _, valid := Value(previous); !valid && previous != "}" && previous != "]" {
				return 0, nil, errors.New("invalid syntax inside array")
			}

			continue

		}

		if previous != "," && previous != "[" {
			return 0, nil, errors.New("invalid previous token")
		}

		value, valid := Value(seg)

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

func Value(seg string) (any, bool) {
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
		return seg[1 : len(seg)-1], true
	}

	if integer, err := strconv.Atoi(seg); err == nil {
		return integer, true
	}

	if f, err := strconv.ParseFloat(seg, 64); err == nil {
		return f, true
	}

	return nil, false
}

func setValue(mode string, name string, value any, body *any) {
	if mode == "array" {
		(*body).([]any)[len((*body).([]any))-1] = value
	} else {
		(*body).(map[string]any)[name] = value
	}
}

func next_nestType(mode string, token string, name string, body *any) *any {

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
