package fox

import (
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

// Currently not used
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

func corsOrigin(origin string, c *Context, allowedOrigins []string) (string, bool) {
	//if origins is not set, send forbidden
	if corsoptions.Origins == nil {
		c.Status(Status.Forbidden)

		return "", false
	}

	for _, org := range allowedOrigins {
		if origin == org || org == "*" {
			return org, true
		}
	}

	return "", false
}

func corsMethod(method string, c *Context, formattedMethods string, allowedMethods []string) (string, bool) {

	if method == "" {
		return "", true
	}

	//if methods is not set, send method not allowed
	if corsoptions.Methods == nil {
		c.Status(Status.MethodNotAllowed)

		return "", false
	}

	for _, mth := range allowedMethods {
		if method == strings.ToUpper(mth) {
			return formattedMethods, true
		} else if mth == "*" {

			return "*", true
		}
	}

	return "", false
}

func corsHeaders(acrh string, formattedHeaders string, allowedHeaders map[string]bool) (string, bool) {

	if acrh == "" {
		return "", true
	}

	if len(allowedHeaders) == 0 {
		return "", false
	}

	//if wildcard exists or no "Access-Control-Request-Headers" header, return true
	if allowedHeaders["*"] {
		return "*", true
	}

	acrh_a := strings.Split(acrh, ", ")

	for _, name := range acrh_a {
		if !allowedHeaders[name] {
			return "", false
		}
	}

	return formattedHeaders, true
}
