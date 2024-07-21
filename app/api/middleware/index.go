package middleware

const (
	ErrInvalidBasicAuth        = "INVALID_BASIC_AUTH"
	ErrInvalidRequestBody      = "INVALID_REQUEST_BODY"
	ErrInvalidServiceSignature = "INVALID_SERVICE_SIGNATURE"
)

type errorResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
