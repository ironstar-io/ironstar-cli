package api

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fatih/color"
)

// retryHTTPWithExpBackoff retries the given function with exponential backoff up to maxRetries times.
// The function should return a result and an error.
// Example time taken
// Retry 1: 2 seconds
// Retry 2: 4 seconds
// Retry 3: 8 seconds
// Retry 4: 16 seconds
// Retry 5: 32 seconds
// Retry 6: 64 seconds
// Total: 126 seconds
func retryHTTPWithExpBackoff(fn func() (*RawResponse, error), maxRetries int) (*RawResponse, error) {
	var lastErr error
	for i := 0; i < maxRetries+1; i++ {
		// Calculate the backoff time
		if i > 0 {
			backoff := time.Duration(2<<uint(i-1)) * time.Second

			color.Red("[%s] Error occurred on retry attempt %d: %s. Retrying in %v...", time.Now().UTC().Format(time.RFC3339), i, lastErr.Error(), backoff)

			time.Sleep(backoff)
		}

		result, err := fn()
		if err == nil {
			return result, nil
		}
		if result != nil && result.StatusCode >= 400 && result.StatusCode < 500 {
			// Return the error for 4xx status codes without retry
			return result, nil
		}

		lastErr = err
	}

	// Return the last error
	return nil, lastErr
}

func getIPFromURL(urlString string) (string, error) {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	host := parsedUrl.Hostname()
	port := parsedUrl.Port()

	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	ip := ips[0].String()
	return fmt.Sprintf("%s:%s", ip, port), nil
}

func getClientIPAddress() (string, error) {
	resp, err := http.Get("https://api.ironstar.io/cdn-cgi/trace")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP status code %d", resp.StatusCode)
	}

	var ip string
	contentType := resp.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/plain") {
		buf := make([]byte, 1024)
		n, _ := resp.Body.Read(buf)
		lines := strings.Split(string(buf[:n]), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "ip=") {
				ip = strings.TrimPrefix(line, "ip=")
				break
			}
		}
	} else {
		return "", fmt.Errorf("unexpected Content-Type %s", contentType)
	}

	return ip, nil
}

func debugLogs(url string, retries int, err error) {
	fmt.Println()
	color.Red("Unable to connect to the Ironstar API after %d retries. Exiting...\n", retries)

	color.Red("Time:       %s\n", time.Now().UTC().Format(time.RFC3339))

	clientIP, ciperr := getClientIPAddress()
	if ciperr != nil {
		color.Red("Client IP:  %s\n", ciperr.Error())
	} else {
		color.Red("Client IP:  %s\n", clientIP)
	}

	remoteIP, riperr := getIPFromURL(url)
	if riperr != nil {
		color.Red("Remote IP:  %s\n", riperr.Error())
	} else {
		color.Red("Remote IP:  %s\n", remoteIP)
	}
	color.Red("Remote URL: %s\n", url)

	color.Red("Error:      %s\n", err.Error())
}
