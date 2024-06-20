package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ironstar-io/ironstar-cli/internal/errs"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/dustin/go-humanize"
)

type Request struct {
	RunTokenRefresh  bool
	Credentials      types.Keylink
	Method           string
	Path             string
	URL              string
	MapStringPayload map[string]interface{}
	BytePayload      []byte
	Retries          int
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
	Total        uint64
	FriendlyName string
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
	fmt.Printf("\rDownloading %s... %s", wc.FriendlyName, strings.Repeat(" ", 12))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading %s... %s", wc.FriendlyName, humanize.IBytes(wc.Total))
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
	err := r.BuildBytePayload()
	if err != nil {
		return nil, err
	}

	r.URL = GetNankaiBaseURL() + r.Path

	res, err := retryHTTPWithExpBackoff(
		func() (*RawResponse, error) {
			return r.HTTPSend()
		}, r.Retries)
	if err != nil {
		debugLogs(r.URL, r.Retries, err)

		return nil, errors.New(errs.IronstarAPIConnectionErrorMsg)
	}

	return res, nil
}

// ArimaSend - Make a HTTP request to the Ironstar upload/download API (ARIMA)
func (r *Request) ArimaSend() (*RawResponse, error) {
	err := r.BuildBytePayload()
	if err != nil {
		return nil, err
	}

	r.URL = GetArimaBaseURL() + r.Path

	res, err := retryHTTPWithExpBackoff(
		func() (*RawResponse, error) {
			return r.HTTPSend()
		}, r.Retries)
	if err != nil {
		debugLogs(r.URL, r.Retries, err)

		return nil, errors.New(errs.IronstarAPIConnectionErrorMsg)
	}

	return res, nil
}

// ArimaSend - Make a HTTP request to the Ironstar upload/download API (ARIMA)
func (r *Request) ArimaDownload(filepath, friendlyName string) (*RawResponse, error) {
	err := r.BuildBytePayload()
	if err != nil {
		return nil, err
	}

	r.URL = GetArimaBaseURL() + r.Path

	res, err := retryHTTPWithExpBackoff(
		func() (*RawResponse, error) {
			return r.HTTPSDownload(filepath, friendlyName)
		}, r.Retries)
	if err != nil {
		debugLogs(r.URL, r.Retries, err)

		return nil, errors.New(errs.IronstarAPIConnectionErrorMsg)
	}

	return res, nil
}

func (r *Request) HTTPSDownload(filepath, friendlyName string) (*RawResponse, error) {
	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.BytePayload))
	if err != nil {
		out.Close()
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
	if err != nil {
		out.Close()
		return nil, err
	}

	ir := &RawResponse{
		StatusCode: resp.StatusCode,
		CallMethod: r.Method,
		CallURL:    r.URL,
		Header:     resp.Header,
	}

	if resp.StatusCode > 399 {
		if resp != nil && resp.Body != nil {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			ir.Body = body
		}
		defer resp.Body.Close()

		// Close the file without defer so it can happen before Rename()
		out.Close()

		if err = os.Remove(filepath + ".tmp"); err != nil {
			return nil, err
		}

		return ir, nil
	}

	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{
		FriendlyName: friendlyName,
	}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return nil, err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print(" - Complete!\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return nil, err
	}

	// chmod 400 to ensure they're read-only for the current user
	err = os.Chmod(filepath, 0400)
	if err != nil {
		return nil, err
	}

	return ir, nil
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
		Body:       bodyBytes,
		StatusCode: resp.StatusCode,
		CallMethod: r.Method,
		CallURL:    r.URL,
		Header:     resp.Header,
	}

	defer resp.Body.Close()

	return ir, nil
}
