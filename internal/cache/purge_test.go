package cache

import (
	"testing"

	"github.com/ironstar-io/ironstar-cli/internal/types"
)

func TestNormalizeCachePurgeURL(t *testing.T) {
	t.Run("no URL preserves full environment purge", func(t *testing.T) {
		purgeURL, err := normalizeCachePurgeURL(nil)
		if err != nil {
			t.Fatal(err)
		}
		if purgeURL != "" {
			t.Fatalf("expected no URL, got %q", purgeURL)
		}
	})

	t.Run("URL is trimmed and preserved", func(t *testing.T) {
		purgeURL, err := normalizeCachePurgeURL([]string{" https://www.example.com/one "})
		if err != nil {
			t.Fatal(err)
		}
		if purgeURL != "https://www.example.com/one" {
			t.Fatalf("unexpected URL %q", purgeURL)
		}
	})

	t.Run("multiple URLs are rejected", func(t *testing.T) {
		if _, err := normalizeCachePurgeURL([]string{"https://www.example.com/one", "https://www.example.com/two"}); err == nil {
			t.Fatal("expected multiple URLs to be rejected")
		}
	})

	t.Run("explicitly empty URL cannot become full purge", func(t *testing.T) {
		if _, err := normalizeCachePurgeURL([]string{"  "}); err == nil {
			t.Fatal("expected an empty URL error")
		}
	})

	for _, value := range []string{
		"http://www.example.com/one",
		"www.example.com/one",
		"https:///missing-host",
		"https://user@example.com/private",
		"https://www.example.com/page#fragment",
	} {
		t.Run("rejects "+value, func(t *testing.T) {
			if _, err := normalizeCachePurgeURL([]string{value}); err == nil {
				t.Fatalf("expected %q to be rejected", value)
			}
		})
	}
}

func TestMakeCacheInvalidationRequest(t *testing.T) {
	t.Run("full purge is hard environment invalidation", func(t *testing.T) {
		request := makeCacheInvalidationRequest("")
		if request.Kind != types.CacheInvalidationKindEnvironment || request.InvalidationType != types.CacheInvalidationTypeHard || request.URL != "" {
			t.Fatalf("unexpected request %#v", request)
		}
	})

	t.Run("selective purge is a soft URL invalidation", func(t *testing.T) {
		request := makeCacheInvalidationRequest("https://www.example.com/one")
		if request.Kind != types.CacheInvalidationKindURL || request.InvalidationType != types.CacheInvalidationTypeSoft {
			t.Fatalf("unexpected request %#v", request)
		}
		if request.URL != "https://www.example.com/one" {
			t.Fatalf("unexpected URL %q", request.URL)
		}
	})
}
