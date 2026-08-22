package cache

import "testing"

func TestNormalizeCachePurgeURLs(t *testing.T) {
	t.Run("no URL preserves full environment purge", func(t *testing.T) {
		urls, err := normalizeCachePurgeURLs(nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(urls) != 0 {
			t.Fatalf("expected no URLs, got %#v", urls)
		}
	})

	t.Run("URLs are trimmed and preserved", func(t *testing.T) {
		urls, err := normalizeCachePurgeURLs([]string{
			" https://www.example.com/one ",
			"https://www.example.com/two",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(urls) != 2 || urls[0] != "https://www.example.com/one" || urls[1] != "https://www.example.com/two" {
			t.Fatalf("unexpected URLs %#v", urls)
		}
	})

	t.Run("explicitly empty URL cannot become full purge", func(t *testing.T) {
		if _, err := normalizeCachePurgeURLs([]string{"  "}); err == nil {
			t.Fatal("expected an empty URL error")
		}
	})
}
