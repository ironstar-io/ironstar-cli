package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ironstar-io/ironstar-cli/internal/types"
)

func TestNewTransferHTTPClientHasNoTotalTimeout(t *testing.T) {
	client := newTransferHTTPClient()

	if client.Timeout != 0 {
		t.Fatalf("expected no total Client.Timeout, got %s", client.Timeout)
	}

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", client.Transport)
	}
	if transport.DialContext == nil {
		t.Fatal("expected DialContext to be set")
	}
	if transport.TLSHandshakeTimeout != transferTLSHandshakeTimeout {
		t.Fatalf("expected TLS handshake timeout %s, got %s", transferTLSHandshakeTimeout, transport.TLSHandshakeTimeout)
	}
	if transport.ResponseHeaderTimeout != transferResponseHeaderTimeout {
		t.Fatalf("expected response header timeout %s, got %s", transferResponseHeaderTimeout, transport.ResponseHeaderTimeout)
	}
}

func TestTransferIdleTimeout(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		t.Setenv("IRONSTAR_TRANSFER_IDLE_TIMEOUT", "")
		if got := transferIdleTimeout(); got != defaultTransferIdleTimeout {
			t.Fatalf("expected default %s, got %s", defaultTransferIdleTimeout, got)
		}
	})

	t.Run("override", func(t *testing.T) {
		t.Setenv("IRONSTAR_TRANSFER_IDLE_TIMEOUT", "45s")
		if got := transferIdleTimeout(); got != 45*time.Second {
			t.Fatalf("expected 45s, got %s", got)
		}
	})

	t.Run("invalid falls back to default", func(t *testing.T) {
		t.Setenv("IRONSTAR_TRANSFER_IDLE_TIMEOUT", "not-a-duration")
		if got := transferIdleTimeout(); got != defaultTransferIdleTimeout {
			t.Fatalf("expected default %s, got %s", defaultTransferIdleTimeout, got)
		}
	})

	t.Run("non-positive falls back to default", func(t *testing.T) {
		t.Setenv("IRONSTAR_TRANSFER_IDLE_TIMEOUT", "0s")
		if got := transferIdleTimeout(); got != defaultTransferIdleTimeout {
			t.Fatalf("expected default %s, got %s", defaultTransferIdleTimeout, got)
		}
	})
}

func TestStallGuardTripsWhenNoProgress(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	pr, pw := io.Pipe()
	t.Cleanup(func() { pw.Close() })

	guard := newStallGuard(cancel, pr, 50*time.Millisecond, -1)

	// Read blocks on the empty pipe, so the watchdog sees no progress and trips.
	go guard.Read(make([]byte, 8))

	select {
	case <-ctx.Done():
		cause := context.Cause(ctx)
		if cause == nil || !strings.Contains(cause.Error(), "stalled") {
			t.Fatalf("expected stall cause, got %v", cause)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("watchdog did not trip on a stalled transfer")
	}
}

func TestStallGuardSurvivesSlowButSteadyProgress(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	// 10 single-byte reads, 20ms apart: total ~200ms (well past the 60ms idle
	// window) but each gap stays under it, so the watchdog must never trip.
	guard := newStallGuard(cancel, &dripReader{max: 10, delay: 20 * time.Millisecond}, 60*time.Millisecond, -1)

	got, err := io.ReadAll(guard)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}
	if len(got) != 10 {
		t.Fatalf("expected 10 bytes, got %d", len(got))
	}
	if cause := context.Cause(ctx); cause != nil {
		t.Fatalf("expected no cancellation, got %v", cause)
	}
}

func TestStallGuardDisarmsAfterExpectedBytes(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	// Reader keeps yielding forever, but the guard should disarm once `expected`
	// bytes have moved, leaving the post-body wait to ResponseHeaderTimeout.
	guard := newStallGuard(cancel, &dripReader{max: -1}, 40*time.Millisecond, 5)

	buf := make([]byte, 1)
	for i := range 5 {
		if _, err := guard.Read(buf); err != nil {
			t.Fatalf("unexpected read error at %d: %v", i, err)
		}
	}

	time.Sleep(120 * time.Millisecond)
	if cause := context.Cause(ctx); cause != nil {
		t.Fatalf("expected watchdog disarmed after expected bytes, got %v", cause)
	}
}

func TestStreamSendStreamsBodyWithContentLength(t *testing.T) {
	tarball := filepath.Join(t.TempDir(), "package.tar.gz")
	contents := strings.Repeat("ironstar", 256)
	if err := os.WriteFile(tarball, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}

	var (
		gotContentLength int64
		gotFile          string
		gotBranch        string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentLength = r.ContentLength

		file, _, err := r.FormFile("package")
		if err != nil {
			t.Errorf("FormFile: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		b, _ := io.ReadAll(file)
		gotFile = string(b)
		gotBranch = r.FormValue("branch")

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	s := &Stream{
		Method:   http.MethodPost,
		URL:      server.URL,
		FilePath: tarball,
		Credentials: types.Keylink{
			AuthToken: "test-token",
		},
		Payload: map[string]string{"branch": "do-not-deploy"},
	}

	resp, err := s.Send()
	if err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
	if string(resp.Body) != `{"ok":true}` {
		t.Fatalf("unexpected response body %q", string(resp.Body))
	}
	if gotContentLength <= 0 {
		t.Fatalf("expected a Content-Length on the upload, got %d", gotContentLength)
	}
	if gotFile != contents {
		t.Fatalf("uploaded file contents did not round-trip")
	}
	if gotBranch != "do-not-deploy" {
		t.Fatalf("expected branch form field, got %q", gotBranch)
	}
}

// dripReader emits one byte per Read, optionally pausing between reads. A
// negative max yields forever; otherwise it returns io.EOF once max bytes have
// been produced.
type dripReader struct {
	max   int
	sent  int
	delay time.Duration
}

func (d *dripReader) Read(p []byte) (int, error) {
	if d.max >= 0 && d.sent >= d.max {
		return 0, io.EOF
	}
	if d.delay > 0 {
		time.Sleep(d.delay)
	}
	d.sent++
	p[0] = 'x'
	return 1, nil
}
