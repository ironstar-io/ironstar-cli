package services

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/types"
)

func TestInitializeWritesIronignoreNotConfigExclude(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := InitializeIronstarProject(); err != nil {
		t.Fatal(err)
	}

	// config.yml must no longer carry a package.exclude list.
	cfg, err := os.ReadFile(filepath.Join(dir, ".ironstar", "config.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(cfg), "exclude") {
		t.Fatalf("config.yml should not contain package.exclude:\n%s", cfg)
	}

	// .ironignore must exist and carry the smart defaults.
	ig, err := os.ReadFile(filepath.Join(dir, ".ironstar", ".ironignore"))
	if err != nil {
		t.Fatalf(".ironignore was not created: %v", err)
	}
	for _, want := range []string{".git", "node_modules", "*.sql", "sites/*/files"} {
		if !strings.Contains(string(ig), want) {
			t.Fatalf(".ironignore missing expected default %q", want)
		}
	}

	// The generated .ironignore must actually drive packaging via the resolver.
	ex, _, err := resolvePackageExcluder(dir, types.ProjectConfig{}, flags.Accumulator{})
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := ex.Excludes("web/themes/x/node_modules", true); !got {
		t.Fatal("expected node_modules to be excluded via the generated .ironignore")
	}
}
