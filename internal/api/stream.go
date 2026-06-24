package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ironstar-io/ironstar-cli/internal/types"
)

type Stream struct {
	RunTokenRefresh bool
	Credentials     types.Keylink
	Method          string
	URL             string
	FilePath        string
	Payload         map[string]string
}

const IronstarUploadAPIDomain = "https://uploads.ironstar.io"

func GetUploadURL(subHash string) string {
	if os.Getenv("IRONSTAR_USE_ARIMA_UPLOAD") != "" {
		// Specifically use Arima for uploads
		domain := IronstarUploadAPIDomain
		if override := os.Getenv("IRONSTAR_UPLOAD_DOMAIN"); override != "" {
			domain = override
		}
		return fmt.Sprintf("%s/upload/subscription/%s", domain, subHash)
	}

	// Default - Use Nankai API
	domain := IronstarProductionAPIDomain
	if override := os.Getenv("IRONSTAR_UPLOAD_DOMAIN"); override != "" {
		domain = override
	}
	return fmt.Sprintf("%s/subscription/%s/upload-package", domain, subHash)
}

// Send - Make a HTTP request to the Ironstar API
func (s *Stream) Send() (*RawResponse, error) {
	file, err := os.Open(s.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("package", filepath.Base(s.FilePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	for k, v := range s.Payload {
		if v == "" {
			continue
		}

		f, err := writer.CreateFormField(k)
		if err != nil {
			return nil, err
		}
		_, err = f.Write([]byte(v))
		if err != nil {
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	bodyLen := int64(body.Len())

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	req, err := http.NewRequestWithContext(ctx, s.Method, s.URL, nil)
	if err != nil {
		return nil, err
	}

	guard := newStallGuard(cancel, body, transferIdleTimeout(), bodyLen)
	defer guard.Stop()
	req.Body = io.NopCloser(guard)
	req.ContentLength = bodyLen

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("authorization", "Bearer "+s.Credentials.AuthToken)
	req.Header.Add("user-agent", fmt.Sprintf("ironstar-cli/%s", version))

	client := newTransferHTTPClient()
	resp, err := client.Do(req)

	var bodyBytes []byte
	if resp != nil && resp.Body != nil {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		bodyBytes = body
	}

	if err != nil {
		if cause := context.Cause(ctx); cause != nil && cause != context.Canceled {
			return nil, cause
		}
		return nil, err
	}

	ir := &RawResponse{
		StatusCode: resp.StatusCode,
		Body:       bodyBytes,
		Header:     resp.Header,
	}

	defer resp.Body.Close()

	return ir, nil
}
