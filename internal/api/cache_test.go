package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ironstar-io/ironstar-cli/internal/types"
)

func TestPostEnvironmentCacheInvalidationBuildsFastlyPayload(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want map[string]interface{}
	}{
		{
			name: "entire environment",
			want: map[string]interface{}{
				"kind":              "environment",
				"invalidation_type": "hard",
			},
		},
		{
			name: "single URL",
			url:  "https://www.example.com/articles/one",
			want: map[string]interface{}{
				"kind":              "url",
				"invalidation_type": "soft",
				"url":               "https://www.example.com/articles/one",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalClient := newAPIHTTPClient
			t.Cleanup(func() { newAPIHTTPClient = originalClient })
			t.Setenv("IRONSTAR_API_ADDRESS", "https://api.example.test")

			var gotPayload map[string]interface{}
			newAPIHTTPClient = func() *http.Client {
				return &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					if req.Method != http.MethodPost {
						t.Fatalf("expected POST request, got %s", req.Method)
					}
					if req.URL.String() != "https://api.example.test/subscription/sub/environment/env/cache-invalidation" {
						t.Fatalf("unexpected request URL %s", req.URL.String())
					}
					if err := json.NewDecoder(req.Body).Decode(&gotPayload); err != nil {
						t.Fatal(err)
					}

					return &http.Response{
						StatusCode: http.StatusCreated,
						Header:     make(http.Header),
						Body:       io.NopCloser(strings.NewReader(`{"name":"purge-test","status":"PENDING"}`)),
					}, nil
				})}
			}

			_, err := PostEnvironmentCacheInvalidation(types.Keylink{}, "text", "sub", "env", tt.url)
			if err != nil {
				t.Fatal(err)
			}

			if len(gotPayload) != len(tt.want) {
				t.Fatalf("expected payload %#v, got %#v", tt.want, gotPayload)
			}
			for key, wantValue := range tt.want {
				if gotPayload[key] != wantValue {
					t.Errorf("expected %s=%q, got %q", key, wantValue, gotPayload[key])
				}
			}
			if _, found := gotPayload["objects"]; found {
				t.Fatal("payload must not use the retired objects field")
			}
		})
	}
}
