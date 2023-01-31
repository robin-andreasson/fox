package fox

import "strings"

type CorsOptions struct {
	Origins     []string
	Methods     []string
	Headers     []string
	Credentials bool

	_formattedMethods string
	_formattedHeaders string
}

var corsoptions = CorsOptions{}

func CORS(options CorsOptions) {
	corsoptions = options

	if options.Methods != nil {
		corsoptions._formattedMethods = formatWithDelimiter(corsoptions.Methods, ", ")
	}

	if options.Headers != nil {
		corsoptions._formattedHeaders = formatWithDelimiter(corsoptions.Headers, ", ")
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
