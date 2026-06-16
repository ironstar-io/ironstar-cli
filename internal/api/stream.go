package api

import (
	"bytes"
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

	req, err := http.NewRequest(s.Method, s.URL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("authorization", "Bearer "+s.Credentials.AuthToken)

	client := newAPIHTTPClient()
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
