package api

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestNewAPIHTTPClientClonesDefaultTransportWithProxyAndTLSSettings(t *testing.T) {
	originalInsecureSkipVerify := flags.Acc.InsecureSkipVerify
	flags.Acc.InsecureSkipVerify = true
	t.Cleanup(func() {
		flags.Acc.InsecureSkipVerify = originalInsecureSkipVerify
	})

	t.Setenv("HTTP_PROXY", "")
	t.Setenv("http_proxy", "")
	t.Setenv("HTTPS_PROXY", "http://192.168.105.101:3128")
	t.Setenv("https_proxy", "")
	t.Setenv("ALL_PROXY", "")
	t.Setenv("all_proxy", "")
	t.Setenv("NO_PROXY", "")
	t.Setenv("no_proxy", "")

	client := newAPIHTTPClient()

	if client.Timeout != apiHTTPClientTimeout {
		t.Fatalf("expected timeout %s, got %s", apiHTTPClientTimeout, client.Timeout)
	}

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", client.Transport)
	}

	if transport == http.DefaultTransport {
		t.Fatal("expected cloned transport, got http.DefaultTransport")
	}

	if transport.Proxy == nil {
		t.Fatal("expected proxy function from default transport")
	}

	req, err := http.NewRequest(http.MethodGet, "https://stage.api.ironstar.io/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}

	proxyURL, err := transport.Proxy(req)
	if err != nil {
		t.Fatal(err)
	}
	if proxyURL == nil || proxyURL.String() != "http://192.168.105.101:3128" {
		t.Fatalf("expected HTTPS proxy, got %v", proxyURL)
	}

	if transport.TLSClientConfig == nil {
		t.Fatal("expected TLSClientConfig to be set")
	}
	if !transport.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("expected InsecureSkipVerify to follow flags.Acc")
	}
}

func TestAPIRequestPathsUseSharedHTTPClient(t *testing.T) {
	originalNewAPIHTTPClient := newAPIHTTPClient
	t.Cleanup(func() {
		newAPIHTTPClient = originalNewAPIHTTPClient
	})

	var requestedURLs []string
	newAPIHTTPClient = func() *http.Client {
		return &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				requestedURLs = append(requestedURLs, req.URL.String())

				switch req.URL.String() {
				case "https://api.example.test/widgets":
					return testResponse(http.StatusCreated, "application/json", `{"ok":true}`), nil
				case "https://uploads.example.test/download":
					return testResponse(http.StatusOK, "application/octet-stream", "downloaded"), nil
				case "https://api.ironstar.io/cdn-cgi/trace":
					return testResponse(http.StatusOK, "text/plain", "ip=203.0.113.10\n"), nil
				default:
					t.Fatalf("unexpected request URL %s", req.URL.String())
					return nil, nil
				}
			}),
		}
	}

	sendReq := &Request{
		Method:      http.MethodPost,
		URL:         "https://api.example.test/widgets",
		BytePayload: []byte(`{"name":"test"}`),
	}
	sendResp, err := sendReq.HTTPSend()
	if err != nil {
		t.Fatal(err)
	}
	if sendResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected HTTPSend status %d, got %d", http.StatusCreated, sendResp.StatusCode)
	}
	if string(sendResp.Body) != `{"ok":true}` {
		t.Fatalf("expected HTTPSend body, got %q", string(sendResp.Body))
	}

	downloadPath := filepath.Join(t.TempDir(), "download.txt")
	downloadReq := &Request{
		Method: http.MethodGet,
		URL:    "https://uploads.example.test/download",
	}
	downloadResp, err := downloadReq.HTTPSDownload(downloadPath, "download.txt")
	if err != nil {
		t.Fatal(err)
	}
	if downloadResp.StatusCode != http.StatusOK {
		t.Fatalf("expected HTTPSDownload status %d, got %d", http.StatusOK, downloadResp.StatusCode)
	}

	downloaded, err := os.ReadFile(downloadPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(downloaded) != "downloaded" {
		t.Fatalf("expected downloaded file contents, got %q", string(downloaded))
	}

	ip, err := getClientIPAddress()
	if err != nil {
		t.Fatal(err)
	}
	if ip != "203.0.113.10" {
		t.Fatalf("expected client IP, got %q", ip)
	}

	expectedURLs := []string{
		"https://api.example.test/widgets",
		"https://uploads.example.test/download",
		"https://api.ironstar.io/cdn-cgi/trace",
	}
	if strings.Join(requestedURLs, "\n") != strings.Join(expectedURLs, "\n") {
		t.Fatalf("expected requested URLs %v, got %v", expectedURLs, requestedURLs)
	}
}

func testResponse(statusCode int, contentType, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header: http.Header{
			"Content-Type": []string{contentType},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
}
