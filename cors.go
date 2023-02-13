package fox

type CorsOptions struct {
	Origins     []string
	Methods     []string
	Headers     []string
	Credentials bool
}

type formattedCorsOptions struct {
	credentials bool

	formattedMethods string
	formattedHeaders string

	mappedOrigins map[string]bool
	mappedMethods map[string]bool
	mappedHeaders map[string]bool
}

var corsoptions = formattedCorsOptions{}

/*
Set Cross-Origin Resource Sharing headers and handle OPTIONS requests
*/
func CORS(options CorsOptions) {

	corsoptions.mappedOrigins = make(map[string]bool)
	corsoptions.mappedMethods = make(map[string]bool)
	corsoptions.mappedHeaders = make(map[string]bool)

	if options.Origins != nil {
		for _, o := range options.Origins {
			corsoptions.mappedOrigins[o] = true
		}
	}

	if options.Methods != nil {
		corsoptions.formattedMethods = formatWithDelimiter(options.Methods, ", ", "*")

		for _, m := range options.Methods {
			corsoptions.mappedMethods[m] = true
		}
	}

	if options.Headers != nil {
		corsoptions.formattedHeaders = formatWithDelimiter(options.Headers, ", ", "*")

		for _, h := range options.Headers {
			corsoptions.mappedHeaders[h] = true
		}
	}

	if options.Credentials {
		corsoptions.credentials = true
	}
}

/*
Checks cors headers and returns the appropiate status code (zero if nothing should happen)
*/
func handleCors(c *Context) int {

	sec_fetch_site := c.Headers["Sec-Fetch-Site"]

	if sec_fetch_site != "cross-site" && sec_fetch_site != "same-site" {
		return 0
	}

	//ORIGIN
	origin_h := c.Headers["Origin"]

	set_acao := c.setHeaders["Access-Control-Allow-Origin"]
	mappedOrigins := corsoptions.mappedOrigins

	if len(set_acao) == 1 {
		origin_h = set_acao[0]
		mappedOrigins = map[string]bool{set_acao[0]: true}
	}

	origin, isAllowedOrigin := validateCors(origin_h, mappedOrigins, origin_h, true)

	if !isAllowedOrigin {
		return Status.Forbidden
	}

	if origin != "" {
		c.SetHeader("Access-Control-Allow-Origin", origin)
	}

	if corsoptions.credentials {
		c.SetHeader("Access-Control-Allow-Credentials", "true")
	}
	//ORIGIN

	acrm := c.Headers["Access-Control-Request-Method"]
	acrh := c.Headers["Access-Control-Request-Headers"]

	if (acrm == "" && acrh == "") || (c.Method != "OPTIONS") {
		return 0
	}

	//METHODS
	set_acam := c.setHeaders["Access-Control-Allow-Methods"]
	mappedMethods := corsoptions.mappedMethods
	formattedMethods := corsoptions.formattedMethods

	if len(set_acam) == 1 {
		mappedMethods = splitComma(set_acam[0])
		formattedMethods = set_acam[0]
	}

	acam, isAllowedMethod := validateCors(acrm, mappedMethods, formattedMethods, true)

	if acam != "" {
		c.SetHeader("Access-Control-Allow-Methods", acam)
	}

	if !isAllowedMethod {
		return Status.MethodNotAllowed
	}
	//METHODS

	//HEADERS
	set_acah := c.setHeaders["Access-Control-Allow-Headers"]
	mappedHeaders := corsoptions.mappedHeaders
	formattedHeaders := corsoptions.formattedHeaders

	if len(set_acah) == 1 {
		mappedHeaders = splitComma(set_acah[0])
		formattedHeaders = set_acah[0]
	}

	acah, isAllowedHeaders := validateCors(acrh, mappedHeaders, formattedHeaders, false)

	if acah != "" {
		c.SetHeader("Access-Control-Allow-Headers", acah)
	}

	if !isAllowedHeaders {
		return Status.Forbidden
	}
	//HEADERS

	return Status.Ok
}

func validateCors(target string, allowedTargets map[string]bool, formattedTargets string, isNotHeaders bool) (string, bool) {

	if target == "" {
		return formattedTargets, true
	}

	if allowedTargets == nil {
		return "", false
	}

	if allowedTargets["*"] {
		return "*", true
	}

	if isNotHeaders {

		if allowedTargets[target] {
			return formattedTargets, true
		}

	} else {

		target_map := splitComma(target)

		for header := range target_map {
			if !allowedTargets[header] {
				return formattedTargets, false
			}
		}

		return formattedTargets, true
	}

	return formattedTargets, false
}
