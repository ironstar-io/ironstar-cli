package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
)

const (
	transferDialTimeout           = 30 * time.Second
	transferKeepAlive             = 30 * time.Second
	transferTLSHandshakeTimeout   = 15 * time.Second
	transferResponseHeaderTimeout = 5 * time.Minute
	defaultTransferIdleTimeout    = 2 * time.Minute
)

// newTransferHTTPClient builds the client used for package uploads and downloads.
// It deliberately has no http.Client.Timeout: that field is a total-request
// deadline (dial + TLS + sending the whole body + reading the response), so on a
// slow link a transfer larger than the budget is killed even while it is making
// steady progress. Liveness here is instead enforced by granular transport
// timeouts plus a per-read idle watchdog (stallGuard), so a transfer runs to
// completion as long as bytes keep moving and fails fast once they stop.
var newTransferHTTPClient = func() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: flags.Acc.InsecureSkipVerify,
	}
	transport.DialContext = (&net.Dialer{
		Timeout:   transferDialTimeout,
		KeepAlive: transferKeepAlive,
	}).DialContext
	transport.TLSHandshakeTimeout = transferTLSHandshakeTimeout
	// Counts only from after the request body is fully written, so it bounds
	// server-side processing without truncating a slow upload.
	transport.ResponseHeaderTimeout = transferResponseHeaderTimeout

	return &http.Client{Transport: transport}
}

// transferIdleTimeout is the window a transfer may make no progress before the
// stall watchdog aborts it. Overridable for slow or flaky links.
func transferIdleTimeout() time.Duration {
	if v := os.Getenv("IRONSTAR_TRANSFER_IDLE_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d
		}
	}

	return defaultTransferIdleTimeout
}

// stallGuard wraps the stream the HTTP transport pulls from (the request body on
// an upload, the response body on a download). Forward progress resets an idle
// timer; if the stream transfers nothing for idleTimeout the request context is
// cancelled, aborting the in-flight transfer. Once `expected` bytes have moved
// the watchdog disarms, leaving any post-body wait to ResponseHeaderTimeout.
type stallGuard struct {
	r           io.Reader
	idleTimeout time.Duration
	expected    int64 // total bytes to move; <= 0 means rely on EOF
	cancel      context.CancelCauseFunc

	mu    sync.Mutex
	timer *time.Timer
	read  int64
	done  bool
}

func newStallGuard(cancel context.CancelCauseFunc, r io.Reader, idleTimeout time.Duration, expected int64) *stallGuard {
	return &stallGuard{
		r:           r,
		idleTimeout: idleTimeout,
		expected:    expected,
		cancel:      cancel,
	}
}

func (sg *stallGuard) Read(p []byte) (int, error) {
	sg.arm()

	n, err := sg.r.Read(p)
	if n > 0 {
		sg.progress(int64(n))
	}
	if err != nil {
		sg.Stop()
	}

	return n, err
}

// arm starts the watchdog on the first read rather than at construction, so the
// idle window covers data transfer only — not dial, TLS, or (for downloads) the
// wait between sending the request and the first response byte.
func (sg *stallGuard) arm() {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	if sg.timer == nil && !sg.done {
		sg.timer = time.AfterFunc(sg.idleTimeout, sg.trip)
	}
}

func (sg *stallGuard) progress(n int64) {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	sg.read += n
	if sg.done {
		return
	}

	if sg.expected > 0 && sg.read >= sg.expected {
		sg.done = true
		if sg.timer != nil {
			sg.timer.Stop()
		}
		return
	}

	if sg.timer != nil {
		sg.timer.Reset(sg.idleTimeout)
	}
}

// Stop disarms the watchdog. Safe to call repeatedly; the caller defers it to
// guarantee the timer cannot fire after the transfer has returned.
func (sg *stallGuard) Stop() {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	sg.done = true
	if sg.timer != nil {
		sg.timer.Stop()
	}
}

func (sg *stallGuard) trip() {
	sg.cancel(fmt.Errorf("transfer stalled: no data transferred for %s", sg.idleTimeout))
}
