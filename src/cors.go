package fox

import (
	"regexp"
	"strings"
)

type CorsOptions struct {
	Origins     []string
	Methods     []string
	Headers     []string
	Credentials bool

	_formattedMethods string
	_formattedHeaders string
	_mappedHeaders    map[string]bool
}

var cors_safelist_rh = map[string]string{
	"accept-language":  "[0-9a-zA-Z]+",
	"content-language": "[0-9a-zA-Z]+",
	"accept":           `^[^"\(\):<>?@\[\\\]{}]+$`,
	"content-type":     `^[^"\(\):<>?@\[\\\]{}]+$`,
}

var corsoptions = CorsOptions{}

func CORS(options CorsOptions) {
	corsoptions = options

	corsoptions._mappedHeaders = make(map[string]bool)

	if options.Methods != nil {
		corsoptions._formattedMethods = formatWithDelimiter(corsoptions.Methods, ", ", "*")
	}

	if options.Headers != nil {
		corsoptions._formattedHeaders = formatWithDelimiter(corsoptions.Headers, ", ", "*")

		for _, h := range options.Headers {
			corsoptions._mappedHeaders[h] = true
		}
	}
}

func corsOrigin(origin string, allowedOrigins []string) (string, bool) {

	for _, org := range allowedOrigins {
		if origin == org || org == "*" {
			return org, true
		}
	}

	return "", false
}

func corsMethod(method string, formattedMethods string, allowedMethods []string) (string, bool) {

	for _, mth := range allowedMethods {
		if method == strings.ToUpper(mth) {
			return formattedMethods, true
		} else if mth == "*" {

			return "*", true
		}
	}

	return "", false
}

func corsHeaders(acrh string, allowedHeaders map[string]bool) bool {

	//if wildcard, return true
	if allowedHeaders["*"] {
		return true
	}

	acrh_a := strings.Split(acrh, ", ")

	for _, name := range acrh_a {
		//if header value length is above 128, return false
		if len(value) > 128 {
			return false
		}

		rex_string := cors_safelist_rh[name]

		//if header isn't inside safelist, check if header is inside set access-control-allow-headers
		if rex_string == "" {
			if !allowedHeaders[name] {
				return false
			}

			continue
		}

		//if header is inside safelist, check value validity
		rex := regexp.MustCompile(rex_string)

		if rex.FindString(value) == "" {
			return false
		}
	}

	return true
}
