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

		header := urldecode(segments[0])
		value := urldecode(segments[1])

		headers[header] = value
	}

	return route[0], route[1], headers
}
