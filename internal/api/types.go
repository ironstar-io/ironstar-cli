package api

import "net/http"

type APIError struct {
	StatusCode    int    `json:"status_code,omitempty"`
	Result        string `json:"result,omitempty"`
	CallURL       string `json:"call_url,omitempty"`
	CallMethod    string `json:"call_method,omitempty"`
	IronstarCode  string `json:"ironstar_code,omitempty"`
	CorrelationId string `json:"correlation_id,omitempty"`
	Message       string `json:"message,omitempty"`
}

type RawResponse struct {
	StatusCode int         `json:"status_code,omitempty"`
	CallURL    string      `json:"call_url,omitempty"`
	CallMethod string      `json:"call_method,omitempty"`
	Header     http.Header `json:"header,omitempty"`
	Body       []byte      `json:"body,omitempty"`
}

type FailureBody struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
