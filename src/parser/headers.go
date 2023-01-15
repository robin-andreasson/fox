package parser

import (
	"strings"
)

func Headers(raw string) (string, string, map[string]string) {

	raw_headers := strings.Split(raw, "\r\n")

	headers := map[string]string{}

	route := strings.Split(raw_headers[0], " ")

	headers["Version"] = route[2]

	for i := 1; i < len(raw_headers); i++ {
		segments := strings.Split(raw_headers[i], ": ")

		headers[segments[0]] = segments[1]
	}

	return route[0], route[1], headers
}
