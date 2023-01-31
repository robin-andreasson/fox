package fox

import "time"

func formatTime(v time.Duration) string {
	return time.Now().Add(time.Millisecond * v).Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

func formatWithDelimiter(arr []string, delimiter string, ignore string) string {
	formatted := ""

	for i := 1; i < len(arr); i++ {

		formatted += delimiter + arr[i]
	}

	return arr[0] + formatted
}
