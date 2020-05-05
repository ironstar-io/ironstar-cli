package api

type APIError struct {
	StatusCode   int
	CallURL      string
	CallMethod   string
	IronstarCode string
	Message      string
}

type RawResponse struct {
	StatusCode int
	CallURL    string
	CallMethod string
	Body       []byte
}

type FailureBody struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
