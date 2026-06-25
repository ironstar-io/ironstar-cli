package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/types"
)

func writeIronignore(t *testing.T, root, body string) {
	t.Helper()
	dir := filepath.Join(root, ".ironstar")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".ironignore"), []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestNewGitIgnoreExcluderFollowsGitSemantics(t *testing.T) {
	root := t.TempDir()
	writeIronignore(t, root, "node_modules\n*.sql\nweb/sites/default/files\n")

	ex, err := newGitIgnoreExcluder(filepath.Join(root, ".ironstar", ".ironignore"), nil)
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		path  string
		isDir bool
		want  bool
	}{
		{"node_modules", true, true},                  // slash-less, root
		{"web/themes/x/node_modules", true, true},     // slash-less, nested at any depth
		{"config/db.sql", false, true},                // glob basename at depth
		{"web/sites/default/files", true, true},       // anchored, exact
		{"deep/web/sites/default/files", true, false}, // anchored must NOT match nested
		{"web/core/index.php", false, false},          // kept
	}
	for _, c := range cases {
		got, err := ex.Excludes(c.path, c.isDir)
		if err != nil {
			t.Fatalf("Excludes(%q): %v", c.path, err)
		}
		if got != c.want {
			t.Fatalf("Excludes(%q, dir=%v) = %v, want %v", c.path, c.isDir, got, c.want)
		}
	}
}

// TestResolvePackageExcluderPrefersIronignore proves precedence: when
// .ironstar/.ironignore exists, config.yml package.exclude is ignored.
func TestResolvePackageExcluderPrefersIronignore(t *testing.T) {
	root := t.TempDir()
	writeIronignore(t, root, "node_modules\n")

	proj := types.ProjectConfig{}
	proj.Package.Exclude = []string{"vendor"} // should be ignored in favour of .ironignore

	ex, _, err := resolvePackageExcluder(root, proj, flags.Accumulator{})
	if err != nil {
		t.Fatal(err)
	}

	if got, _ := ex.Excludes("web/x/node_modules", true); !got {
		t.Fatal("expected node_modules (from .ironignore) to be excluded")
	}
	if got, _ := ex.Excludes("vendor", true); got {
		t.Fatal("expected config.yml package.exclude (vendor) to be IGNORED when .ironignore is present")
	}
}

// TestResolvePackageExcluderFallsBackToConfig exercises the legacy path plus the
// slash-less-matches-at-depth improvement.
func TestResolvePackageExcluderFallsBackToConfig(t *testing.T) {
	root := t.TempDir() // no .ironignore

	proj := types.ProjectConfig{}
	proj.Package.Exclude = []string{".git", "vendor"}

	ex, _, err := resolvePackageExcluder(root, proj, flags.Accumulator{})
	if err != nil {
		t.Fatal(err)
	}

	if got, _ := ex.Excludes("web/core/.git", true); !got {
		t.Fatal("expected nested .git to be excluded by the improved legacy matcher")
	}
	if got, _ := ex.Excludes("vendor", true); !got {
		t.Fatal("expected vendor to be excluded by config.yml package.exclude")
	}
}
