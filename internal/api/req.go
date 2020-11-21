package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/internal/services"
	"gitlab.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
)

type Request struct {
	RunTokenRefresh  bool
	Credentials      types.Keylink
	Method           string
	Path             string
	URL              string
	MapStringPayload map[string]interface{}
	BytePayload      []byte
}

const IronstarProductionAPIDomain = "https://api.ironstar.io"
const IronstarArimaProductionAPIDomain = "https://uploads.ironstar.io"

func GetNankaiBaseURL() string {
	ipa := os.Getenv("IRONSTAR_API_ADDRESS")
	if ipa != "" {
		return ipa
	}

	return IronstarProductionAPIDomain
}

func GetArimaBaseURL() string {
	ipa := os.Getenv("IRONSTAR_ARIMA_API_ADDRESS")
	if ipa != "" {
		return ipa
	}

	return IronstarArimaProductionAPIDomain
}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
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

// NankaiSend - Make a HTTP request to the Ironstar API
func (r *Request) NankaiSend() (*RawResponse, error) {
	r.RefreshToken()

	err := r.BuildBytePayload()
	if err != nil {
		return nil, err
	}

	r.URL = GetNankaiBaseURL() + r.Path

	return r.HTTPSend()
}

// ArimaSend - Make a HTTP request to the Ironstar upload/download API (ARIMA)
func (r *Request) ArimaSend() (*RawResponse, error) {
	r.RefreshToken()

	err := r.BuildBytePayload()
	if err != nil {
		return nil, err
	}

	r.URL = GetArimaBaseURL() + r.Path

	return r.HTTPSend()
}

// ArimaSend - Make a HTTP request to the Ironstar upload/download API (ARIMA)
func (r *Request) ArimaDownload(filepath string) error {
	r.RefreshToken()

	err := r.BuildBytePayload()
	if err != nil {
		return err
	}

	r.URL = GetArimaBaseURL() + r.Path

	return r.HTTPSDownload(filepath)
}

func (r *Request) HTTPSDownload(filepath string) error {
	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.BytePayload))
	if err != nil {
		out.Close()
		return err
	}

	req.Header.Add("content-type", "application/json")
	if r.Credentials.AuthToken != "" {
		req.Header.Add("authorization", "Bearer "+r.Credentials.AuthToken)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		out.Close()
		return err
	}

	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}

	fmt.Println("Download Finished")

	return nil
}

func (r *Request) HTTPSend() (*RawResponse, error) {
	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.BytePayload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	if r.Credentials.AuthToken != "" {
		req.Header.Add("authorization", "Bearer "+r.Credentials.AuthToken)
	}

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
		Body:       bodyBytes,
		StatusCode: resp.StatusCode,
		CallMethod: r.Method,
		CallURL:    r.URL,
	}

	defer resp.Body.Close()

	return ir, nil
}

func (r *Request) RefreshToken() {
	//
	// If any leg fails, return silently, user will need to relog manually
	//

	if !r.RunTokenRefresh || r.Credentials == (types.Keylink{}) || r.Credentials.Login == "" || r.Credentials.AuthToken == "" || r.Credentials.Expiry.IsZero() {
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
		MapStringPayload: map[string]interface{}{
			"expiry": time.Now().AddDate(0, 0, 14).UTC().Format(time.RFC3339),
		},
	}

	res, err := newReq.NankaiSend()
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
}
