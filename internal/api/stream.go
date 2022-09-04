package api

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
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
	Path            string
	FilePath        string
	Payload         map[string]string
}

const IronstarUploadAPIDomain = "https://uploads.ironstar.io"

func GetBaseUploadURL() string {
	ipa := os.Getenv("IRONSTAR_ARIMA_API_ADDRESS")
	if ipa != "" {
		return ipa
	}

	return IronstarUploadAPIDomain
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

	url := GetBaseUploadURL() + s.Path
	req, err := http.NewRequest(s.Method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("authorization", "Bearer "+s.Credentials.AuthToken)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)

	var bodyBytes []byte
	if resp != nil && resp.Body != nil {
		body, err := ioutil.ReadAll(resp.Body)
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
