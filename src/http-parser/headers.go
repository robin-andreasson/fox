package parser

import (
	"fmt"
	"strings"
)

func Headers(raw string) map[string]string {

	fmt.Println(strings.Split(raw, "\r\n"))

	return make(map[string]string)
}
