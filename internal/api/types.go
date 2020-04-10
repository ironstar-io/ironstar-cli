package api

type APIError struct {
	StatusCode   int
	IronstarCode string
	Message      string
}

type RawResponse struct {
	StatusCode int
	Body       []byte
}

type FailureBody struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
