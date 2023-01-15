package fox

import "time"

func formatTime(v time.Duration) string {
	return time.Now().Add(time.Millisecond * v).Format("Mon, 02 Jan 2006 15:04:05 GMT")
}
