package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
)

type Request struct {
	RunTokenRefresh  bool
	Credentials      types.Keylink
	Method           string
	Path             string
	MapStringPayload map[string]string
	BytePayload      []byte
}

const IronstarProductionAPIDomain = "https://api.ironstar.io"

func GetBaseURL() string {
	ipa := os.Getenv("IRONSTAR_API_ADDRESS")
	if ipa != "" {
		return ipa
	}

	return IronstarProductionAPIDomain
}

func (r *Request) BuildBytePayload() error {
	if r.MapStringPayload != nil {
		b, err := json.Marshal(r.MapStringPayload)
		if err != nil {
			return err
		}

		r.BytePayload = b
	}

	return nil
}

// Send - Make a HTTP request to the Ironstar API
func (r *Request) Send() (*RawResponse, error) {
	r.RefreshToken()

	err := r.BuildBytePayload()
	if err != nil {
		return nil, err
	}

	url := GetBaseURL() + r.Path
	req, err := http.NewRequest(r.Method, url, bytes.NewBuffer(r.BytePayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+r.Credentials.AuthToken)

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

func (r *Request) RefreshToken() {
	//
	// If any leg fails, return silently, user will need to relog manually
	//

	if r.RunTokenRefresh != true || r.Credentials == (types.Keylink{}) || r.Credentials.Login == "" || r.Credentials.AuthToken == "" || r.Credentials.Expiry.IsZero() {
		return
	}

	// Get the time difference in days
	expDiff := int(math.RoundToEven(r.Credentials.Expiry.Sub(time.Now().UTC()).Hours() / 24))
	// Only refresh if there's less than seven days remaining and not already expired
	if expDiff > 7 || expDiff < 0 {
		return
	}

	newReq := Request{
		RunTokenRefresh: false,
		Credentials: types.Keylink{
			AuthToken: r.Credentials.AuthToken,
		},
		Method: "POST",
		Path:   "/auth/token/refresh",
		MapStringPayload: map[string]string{
			"expiry": time.Now().AddDate(0, 0, 14).UTC().Format(time.RFC3339),
		},
	}

	res, err := newReq.Send()
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		return
	}

	b := &types.AuthResponseBody{}
	err = json.Unmarshal(res.Body, b)
	if err != nil {
		return
	}

	newCreds := types.Keylink{
		Login:     r.Credentials.Login,
		AuthToken: b.IDToken,
		Expiry:    b.Expiry,
	}

	err = services.UpdateCredentialsFile(newCreds)
	if err != nil {
		return
	}

	r.Credentials = newCreds

	color.Yellow("Authentication token has been automatically refreshed")
	fmt.Println()

	return
}
