package api

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/ironstar-io/ironstar-cli/internal/types"
)

func TestDownloadBackupThroughNankai(t *testing.T) {
	for _, tc := range []struct {
		name          string
		apiStatus     int
		body          string
		storageStatus int
		retry         bool
		wantError     bool
		wantTransfers int
	}{
		{"success", 200, `{"downloadURL":"https://storage.example.test/archive?signature=private","size":7}`, 200, false, false, 1},
		{"retry transport failure", 200, `{"downloadURL":"https://storage.example.test/archive?signature=private"}`, 200, true, false, 2},
		{"permission denied", 403, `{"message":"denied"}`, 0, false, true, 0},
		{"invalid JSON", 200, `{`, 0, false, true, 0},
		{"missing URL", 200, `{}`, 0, false, true, 0},
		{"invalid URL", 200, `{"downloadURL":"file:///archive"}`, 0, false, true, 0},
		{"storage denied", 200, `{"downloadURL":"https://storage.example.test/archive?signature=private"}`, 403, false, true, 1},
		{"empty storage response", 200, `{"downloadURL":"https://storage.example.test/archive?signature=private"}`, 204, false, true, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			oldAPI, oldTransfer := newAPIHTTPClient, newTransferHTTPClient
			t.Cleanup(func() { newAPIHTTPClient, newTransferHTTPClient = oldAPI, oldTransfer })
			t.Setenv("IRONSTAR_API_ADDRESS", "https://api.example.test")
			t.Setenv("IRONSTAR_UPLOAD_DOMAIN", "https://unused.example.test")
			newAPIHTTPClient = func() *http.Client {
				return &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					if req.Method != http.MethodGet || req.URL.String() != "https://api.example.test/subscription/sub/environment/live/backups/daily/download-link?component=db+%26+files" {
						t.Fatalf("unexpected API request: %s %s", req.Method, req.URL)
					}
					if req.Header.Get("Authorization") != "Bearer test-token" {
						t.Fatal("missing API authorization")
					}
					return testResponse(tc.apiStatus, "application/json", tc.body), nil
				})}
			}
			transfers := 0
			newTransferHTTPClient = func() *http.Client {
				return &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					transfers++
					if req.Method != http.MethodGet || req.URL.String() != "https://storage.example.test/archive?signature=private" {
						t.Fatal("signed download URL was changed")
					}
					if req.Header.Get("Authorization") != "" {
						t.Fatal("Ironstar token sent to storage")
					}
					if tc.retry && transfers == 1 {
						return nil, errors.New("temporary transfer failure")
					}
					return testResponse(tc.storageStatus, "application/octet-stream", "archive"), nil
				})}
			}
			dest := filepath.Join(t.TempDir(), "backup.tar.gz")
			if tc.wantError {
				if err := os.WriteFile(dest, []byte("existing"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			err := DownloadEnvironmentBackupComponent(types.Keylink{AuthToken: "test-token"}, "", "sub", "live", "daily", dest, types.BackupIterationComponent{Name: "db & files"})
			if (err != nil) != tc.wantError {
				t.Fatalf("unexpected error: %v", err)
			}
			if transfers != tc.wantTransfers {
				t.Fatalf("got %d transfers, want %d", transfers, tc.wantTransfers)
			}
			data, err := os.ReadFile(dest)
			if err != nil {
				t.Fatal(err)
			}
			want := "archive"
			if tc.wantError {
				want = "existing"
			}
			if string(data) != want {
				t.Fatalf("file contains %q, want %q", data, want)
			}
			if _, err := os.Stat(dest + ".tmp"); !os.IsNotExist(err) {
				t.Fatalf("temporary file remains: %v", err)
			}
			if !tc.wantError {
				info, err := os.Stat(dest)
				if err != nil {
					t.Fatal(err)
				}
				if info.Mode().Perm() != 0400 {
					t.Fatalf("unexpected permissions: %v", info.Mode())
				}
			}
		})
	}
}

func TestUploadURLUsesNankai(t *testing.T) {
	t.Setenv("IRONSTAR_API_ADDRESS", "")
	t.Setenv("IRONSTAR_UPLOAD_DOMAIN", "")
	if got := GetUploadURL("sub"); got != "https://api.ironstar.io/subscription/sub/upload-package" {
		t.Fatal(got)
	}
	t.Setenv("IRONSTAR_API_ADDRESS", "https://api.example.test")
	if got := GetUploadURL("sub"); got != "https://api.example.test/subscription/sub/upload-package" {
		t.Fatal(got)
	}
	t.Setenv("IRONSTAR_UPLOAD_DOMAIN", "https://upload.example.test")
	if got := GetUploadURL("sub"); got != "https://upload.example.test/subscription/sub/upload-package" {
		t.Fatal(got)
	}
}
