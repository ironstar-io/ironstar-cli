package cache

import (
	"testing"

	"github.com/ironstar-io/ironstar-cli/internal/types"
)

func TestNormalizeCacheInvalidationURL(t *testing.T) {
	t.Run("no URL selects the environment", func(t *testing.T) {
		invalidationURL, err := normalizeCacheInvalidationURL(nil)
		if err != nil {
			t.Fatal(err)
		}
		if invalidationURL != "" {
			t.Fatalf("expected no URL, got %q", invalidationURL)
		}
	})

	t.Run("URL is trimmed and preserved", func(t *testing.T) {
		invalidationURL, err := normalizeCacheInvalidationURL([]string{" https://www.example.com/one "})
		if err != nil {
			t.Fatal(err)
		}
		if invalidationURL != "https://www.example.com/one" {
			t.Fatalf("unexpected URL %q", invalidationURL)
		}
	})

	t.Run("multiple URLs are rejected", func(t *testing.T) {
		if _, err := normalizeCacheInvalidationURL([]string{"https://www.example.com/one", "https://www.example.com/two"}); err == nil {
			t.Fatal("expected multiple URLs to be rejected")
		}
	})

	t.Run("explicitly empty URL cannot become an environment invalidation", func(t *testing.T) {
		if _, err := normalizeCacheInvalidationURL([]string{"  "}); err == nil {
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
			if _, err := normalizeCacheInvalidationURL([]string{value}); err == nil {
				t.Fatalf("expected %q to be rejected", value)
			}
		})
	}
}

func TestMakeCacheInvalidationRequest(t *testing.T) {
	t.Run("environment invalidation defaults to hard", func(t *testing.T) {
		request, err := makeCacheInvalidationRequest("", "")
		if err != nil {
			t.Fatal(err)
		}
		if request.Kind != types.CacheInvalidationKindEnvironment || request.InvalidationType != types.CacheInvalidationTypeHard || request.URL != "" {
			t.Fatalf("unexpected request %#v", request)
		}
	})

	t.Run("URL invalidation defaults to soft", func(t *testing.T) {
		request, err := makeCacheInvalidationRequest("https://www.example.com/one", "")
		if err != nil {
			t.Fatal(err)
		}
		if request.Kind != types.CacheInvalidationKindURL || request.InvalidationType != types.CacheInvalidationTypeSoft {
			t.Fatalf("unexpected request %#v", request)
		}
		if request.URL != "https://www.example.com/one" {
			t.Fatalf("unexpected URL %q", request.URL)
		}
	})

	t.Run("URL invalidation may be hard", func(t *testing.T) {
		request, err := makeCacheInvalidationRequest("https://www.example.com/one", "HARD")
		if err != nil {
			t.Fatal(err)
		}
		if request.InvalidationType != types.CacheInvalidationTypeHard {
			t.Fatalf("unexpected request %#v", request)
		}
	})

	t.Run("environment invalidation may not be soft", func(t *testing.T) {
		if _, err := makeCacheInvalidationRequest("", "soft"); err == nil {
			t.Fatal("expected soft environment invalidation to be rejected")
		}
	})

	t.Run("unknown invalidation type is rejected", func(t *testing.T) {
		if _, err := makeCacheInvalidationRequest("https://www.example.com/one", "stale"); err == nil {
			t.Fatal("expected unknown invalidation type to be rejected")
		}
	})
}
