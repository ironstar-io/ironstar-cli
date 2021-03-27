package api

import "net/http"

type APIError struct {
	StatusCode    int
	CallURL       string
	CallMethod    string
	IronstarCode  string
	CorrelationId string
	Message       string
}

type RawResponse struct {
	StatusCode int
	CallURL    string
	CallMethod string
	Header     http.Header
	Body       []byte
}

type FailureBody struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
