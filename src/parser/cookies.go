package parser

import (
	b64 "encoding/base64"
	"strings"
)

/*
Parses incoming cookies and maps it into a map with string values
*/
func Cookies(s string) map[string]string {

	cookies := make(map[string]string)

	cookie_seg := strings.Split(s, "; ")

	for _, cookie := range cookie_seg {

		name, value, found := strings.Cut(cookie, "=")

		if !found {
			continue
		}

		if v, err := b64.StdEncoding.DecodeString(value); err == nil {
			value = string(v)
		}

		cookies[name] = value
	}

	return cookies
}
