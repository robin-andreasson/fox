package parser

type mime_patterns struct {
	byte_pattern []byte
	pattern_mask []byte
	ignored      map[byte]bool
	mime         string
}

var extMimes = map[string]string{
	"html": "text/html; charset=utf-8",
	"htm":  "text/html; charset=utf-8",
	"css":  "text/css",
	"js":   "text/javascript",
	"php":  "application/x-httpd-php",
	"xml":  "application/xml",
}

var mimes = []mime_patterns{
	//IMAGE
	{mime: "image/x-icon", byte_pattern: []byte{0x00, 0x00, 0x00, 0x01, 0x00}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "image/x-icon", byte_pattern: []byte{0x00, 0x00, 0x00, 0x02, 0x00}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "image/bmp", byte_pattern: []byte{0x42, 0x4D}, pattern_mask: []byte{0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "image/gif", byte_pattern: []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "image/gif", byte_pattern: []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "image/webp", byte_pattern: []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50, 0x56, 0x50}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "image/png", byte_pattern: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "image/jpeg", byte_pattern: []byte{0xFF, 0xD8, 0xFF}, pattern_mask: []byte{0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	//IMAGE

	//FONT
	{mime: "application/vnd.ms-fontobject", byte_pattern: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x4C, 0x50}, pattern_mask: []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "font/ttf", byte_pattern: []byte{0, 1, 0, 0}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "font/otf", byte_pattern: []byte{0x4F, 0x54, 0x54, 0x4F}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "font/collection", byte_pattern: []byte{0x74, 0x74, 0x63, 0x66}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "font/woff", byte_pattern: []byte{0x77, 0x4F, 0x46, 0x46}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "font/woff2", byte_pattern: []byte{0x77, 0x4F, 0x46, 0x32}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	//FONT

	//AUDIO OR VIDEO
	{mime: "audio/aiff", byte_pattern: []byte{0x46, 0x4F, 0x52, 0x4D, 0, 0, 0, 0, 0x41, 0x49, 0x46, 0x46}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0, 0, 0, 0, 0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "audio/mpeg", byte_pattern: []byte{0x49, 0x44, 0x33}, pattern_mask: []byte{0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "application/ogg", byte_pattern: []byte{0x4F, 0x67, 0x67, 0x53, 0}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "audio/midi", byte_pattern: []byte{0x4D, 0x54, 0x68, 0x64, 0, 0, 0, 0x06}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "video/avi", byte_pattern: []byte{0x52, 0x49, 0x46, 0x46, 0, 0, 0, 0, 0x41, 0x56, 0x49, 0x20}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0, 0, 0, 0, 0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "audio/wave", byte_pattern: []byte{0x52, 0x49, 0x46, 0x46, 0, 0, 0, 0, 0x57, 0x41, 0x56, 0x45}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0, 0, 0, 0, 0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	//AUDIO OR VIDEO

	//ZIP
	{mime: "application/x-gzip", byte_pattern: []byte{0x1F, 0x8B, 0x08}, pattern_mask: []byte{0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "application/zip", byte_pattern: []byte{0x50, 0x4B, 0x03, 0x04}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	{mime: "application/x-gzip", byte_pattern: []byte{0x52, 0x61, 0x72, 0x20, 0x1A, 0x07, 0x00}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	//ZIP

	//PDF
	{mime: "application/pdf", byte_pattern: []byte{0x25, 0x50, 0x44, 0x46, 0x2D}, pattern_mask: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, ignored: map[byte]bool{}},
	//PDF
}

func Mime(input []byte) string {

	for _, mime_pattern := range mimes {

		match := pattern_match(input, mime_pattern.byte_pattern, mime_pattern.pattern_mask, mime_pattern.ignored)

		if match {
			return mime_pattern.mime
		}
	}

	if mp4(input) {
		return "video/mp4"
	}

	return "text/plain"
}

func pattern_match(input []byte, pattern []byte, mask []byte, ignored map[byte]bool) bool {

	if len(pattern) > len(input) {
		return false
	}

	s := 0

	for s < len(input) {
		if !ignored[input[s]] {
			break
		}

		s++
	}

	for i, p := range pattern {
		maskedData := input[s+i] & mask[i]

		if p != maskedData {
			return false
		}
	}

	return true
}

func mp4(input []byte) bool {

	length := len(input)

	if length < 16 {
		return false
	}

	box_size := int(input[3])

	if length < box_size || box_size%4 != 0 {
		return false
	}

	predefinedsignature := string(input[4:8])

	if predefinedsignature != "ftyp" {
		return false
	}

	signature := string(input[8:11])

	if signature == "mp4" {
		return true
	}

	bytes_read := 16

	for bytes_read < box_size {
		signature = string(input[bytes_read : bytes_read+2])

		if signature == "mp4" {
			return true
		}

		bytes_read += 4
	}

	return false
}
