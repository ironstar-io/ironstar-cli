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
	// BytesSent / BodySize report upload progress at the point the response was
	// received, so a server error mid-transfer can be told apart from one that
	// fired only after the whole body was sent.
	BytesSent int64 `json:"bytes_sent,omitempty"`
	BodySize  int64 `json:"body_size,omitempty"`
}

type FailureBody struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
