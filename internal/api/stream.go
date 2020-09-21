package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	// "gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"
)

type Stream struct {
	RunTokenRefresh  bool
	Credentials      types.Keylink
	Method           string
	Path             string
	FilePath         string
	MapStringPayload map[string]interface{}
	BytePayload      []byte
	Ref              string
}

const IronstarUploadAPIDomain = "https://uploads.ironstar.io"

func GetBaseUploadURL() string {
	ipa := os.Getenv("IRONSTAR_UPLOAD_API_ADDRESS")
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
	ref, err := writer.CreateFormField("ref")
	if err != nil {
		return nil, err
	}
	_, err = ref.Write([]byte(s.Ref))
	if err != nil {
		return nil, err
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
	}

	defer resp.Body.Close()

	return ir, nil
}

func (s *Stream) BuildBytePayload() error {
	if s.MapStringPayload != nil {
		b, err := json.Marshal(s.MapStringPayload)
		if err != nil {
			return err
		}

		s.BytePayload = b
	}

	return nil
}
