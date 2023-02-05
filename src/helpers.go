package fox

import (
	"regexp"
	"time"
)

func formatTime(t int) string {
	return time.Now().Add(time.Millisecond * time.Duration(t)).Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

func formatWithDelimiter(arr []string, delimiter string, ignore string) string {

	if len(arr) == 0 {
		return ""
	}

	formatted := ""

	for i := 1; i < len(arr); i++ {
		if arr[i] == ignore {
			continue
		}

		formatted += delimiter + arr[i]
	}

	return arr[0] + formatted
}

// split access control request
func splitComma(target string) map[string]bool {
	rex := regexp.MustCompile(`\s+,\s+|\s+,|,\s+|,`)

	arr := rex.Split(target, -1)

	mappedTarget := map[string]bool{}

	for _, a := range arr {
		mappedTarget[a] = true
	}

	return mappedTarget
}
