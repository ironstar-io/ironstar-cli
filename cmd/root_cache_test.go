package cmd

import "testing"

func TestCacheCommandRoutes(t *testing.T) {
	tests := []struct {
		path   []string
		use    string
		hidden bool
	}{
		{path: []string{"cache"}, use: "cache", hidden: false},
		{path: []string{"cache", "invalidate"}, use: "invalidate [flags]", hidden: false},
		{path: []string{"cache", "list"}, use: "list", hidden: false},
		{path: []string{"cache", "show"}, use: "show [name]", hidden: false},
		{path: []string{"cache", "invalidation"}, use: "invalidation", hidden: true},
		{path: []string{"cache", "invalidation", "create"}, use: "create [flags]", hidden: true},
		{path: []string{"cache", "invalidation", "list"}, use: "list", hidden: true},
		{path: []string{"cache", "invalidation", "show"}, use: "show [name]", hidden: true},
	}

	for _, tt := range tests {
		command, remaining, err := rootCmd.Find(tt.path)
		if err != nil {
			t.Fatalf("find %v: %v", tt.path, err)
		}
		if len(remaining) != 0 {
			t.Fatalf("find %v left arguments %v", tt.path, remaining)
		}
		if command.Use != tt.use {
			t.Errorf("find %v: expected use %q, got %q", tt.path, tt.use, command.Use)
		}
		if command.Hidden != tt.hidden {
			t.Errorf("find %v: expected hidden=%v, got %v", tt.path, tt.hidden, command.Hidden)
		}
	}
}
