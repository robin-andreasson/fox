package fox

type status struct {
	Continue           int
	SwitchingProtocols int
	Processing         int
	EarlyHints         int

	Ok                          int
	Created                     int
	Accepted                    int
	NonAuthoritativeInformation int
	NoContent                   int
	ResetContent                int
	PartialContent              int
	MultiStatus                 int
	AlreadyReported             int
	IMUsed                      int

	MultipleChoices    int
	MovedPermanently   int
	Found              int
	SeeOther           int
	NotModified        int
	UseProxy           int
	Unused             int
	TemeporaryRedirect int
	PermanentRedirect  int

	BadRequest                  int
	Unauthorized                int
	PaymentRequired             int
	Forbidden                   int
	NotFound                    int
	MethodNotAllowed            int
	NotAcceptable               int
	ProxyAuthenticationRequired int
	RequestTimeout              int
	Conflict                    int
	Gone                        int
	LengthRequired              int
	PreconditionFailed          int
	PayloadTooLarge             int
	URITooLong                  int
	UnsupportedMediaType        int
	RangeNotSatisfiable         int
	ExpectationFailed           int
	ImaTeapot                   int
	MisdirectedRequest          int
	UnprocessableEntity         int
	Locked                      int
	FailedDependency            int
	TooEarly                    int
	UpgradeRequired             int
	PreconditionRequired        int
	TooManyRequests             int
	RequestHeaderFieldsToolarge int
	UnavailableForLegalReason   int

	InternalServerError           int
	NotImplemented                int
	BadGateway                    int
	ServiceUnavailable            int
	GatewayTimeout                int
	HTTPVersionNotSupported       int
	VariantAlsoNegotiates         int
	InsufficientStorage           int
	LoopDetected                  int
	NotExtended                   int
	NetworkAuthenticationRequired int
}

var Status status = status{
	100,
	101,
	102,
	103,

	200,
	201,
	202,
	203,
	204,
	205,
	206,
	207,
	208,
	226,

	300,
	301,
	302,
	303,
	304,
	305,
	306,
	307,
	308,

	400,
	401,
	402,
	403,
	404,
	405,
	406,
	407,
	408,
	409,
	410,
	411,
	412,
	413,
	414,
	415,
	416,
	417,
	418,
	421,
	422,
	423,
	424,
	425,
	426,
	428,
	429,
	431,
	451,

	500,
	501,
	502,
	503,
	504,
	505,
	506,
	507,
	508,
	510,
	511,
}
