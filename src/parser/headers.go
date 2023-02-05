package parser

import (
	"strings"
)

func Headers(raw string) (string, string, map[string]string) {

	raw_headers := strings.Split(raw, "\r\n")

	headers := map[string]string{}

	route := strings.Split(raw_headers[0], " ")

	headers["Protocol"] = route[2]

	for i := 1; i < len(raw_headers); i++ {
		header, value, found := strings.Cut(raw_headers[i], ": ")

		if !found {
			continue
		}

		header = strings.Title(strings.ToLower(urldecode(header)))
		value = urldecode(value)

		headers[header] = value
	}

	return route[0], route[1], headers
}
