package parser

import (
	"reflect"
	"strconv"
)

// Splits array at a specific target
func FirstInstance(data []byte, target string) ([]byte, []byte, bool) {

	index := -1
	found := 0
	targetlength := len(target)

	for i := 0; i < len(data); i++ {

		if data[i] == target[found] {
			found++

			if found == targetlength {
				index = i + 1
				break
			}

		} else {
			found = 0
		}
	}

	if index == -1 {
		return nil, nil, false
	}

	return data[0 : index-targetlength], data[index:], true
}

func ExtensionMime(path string) (string, bool) {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {

			mime := extMimes[path[i+1:]]

			if mime == "" {
				return "", false
			}

			return mime, true
		}
	}

	return "", false
}

func urldecode(s string) string {

	ns := ""

	for i := 0; i < len(s); i++ {
		char := s[i]

		if char != '%' || i+3 > len(s) {
			ns += string(char)

			continue
		}

		hex := s[i+1 : i+3]

		decimal, err := strconv.ParseInt(hex, 16, 32)

		if err != nil {
			continue
		}

		i += 2

		ns += string(decimal)
	}

	return ns
}

func unicodedecode(s string) string {

	offset := 0

	for index, char := range s {

		i := index - offset

		if char != '\\' || i+6 > len(s) {
			continue
		}

		if s[i+1] != 'u' {
			continue
		}

		unicode := s[i+1 : i+6]

		code := unicode[1:]

		decimal, err := strconv.ParseInt(code, 16, 32)

		if err != nil {
			continue
		}

		decimal_s := string(decimal)

		offset += 6 - len(decimal_s)

		s = s[0:i] + decimal_s + s[i+6:]
	}

	return s
}

func IsMap(v any) bool {
	if v == nil {
		return false
	}

	return reflect.TypeOf(v).Kind() == reflect.Map
}

func IsArray(v any) bool {
	if v == nil {
		return false
	}

	t := reflect.TypeOf(v).Kind()

	if t != reflect.Array && t != reflect.Slice {
		return false
	}

	return true
}

func isBytes(v any) bool {
	if v == nil {
		return false
	}

	return reflect.TypeOf(v).Elem().Kind() == reflect.Uint8
}

func getNumber(v string) (any, bool) {

	if integer, err := strconv.Atoi(v); err == nil {
		return integer, true
	}

	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return f, true
	}

	return nil, false
}

//first index: starts directly at that index
//last index: starts at that index minus 1
