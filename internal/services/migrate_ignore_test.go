package services

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRemovePackageBlockPreservesOtherKeys(t *testing.T) {
	in := `version: "1.0"
subscription:
  alias: spacestation
package:
  exclude:
  - .git
  - private
tasks:
  deploy:
    - command: drush updb
`
	out, removed := removePackageBlock([]byte(in))
	if !removed {
		t.Fatal("expected the package block to be removed")
	}
	s := string(out)

	for _, gone := range []string{"package:", "exclude:", "- .git", "- private"} {
		if strings.Contains(s, gone) {
			t.Fatalf("expected %q removed, got:\n%s", gone, s)
		}
	}
	for _, keep := range []string{`version: "1.0"`, "subscription:", "alias: spacestation", "tasks:", "drush updb"} {
		if !strings.Contains(s, keep) {
			t.Fatalf("expected %q preserved, got:\n%s", keep, s)
		}
	}
}

func TestRemovePackageBlockNoOpWhenAbsent(t *testing.T) {
	in := "version: \"1.0\"\ntasks:\n  deploy: []\n"
	out, removed := removePackageBlock([]byte(in))
	if removed {
		t.Fatal("expected no removal when there is no package block")
	}
	if string(out) != in {
		t.Fatal("expected bytes unchanged when no package block is present")
	}
}

func TestMigrateToIronignore(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".ironstar"), 0755); err != nil {
		t.Fatal(err)
	}
	cfg := `version: "1.0"
subscription:
  alias: demo
package:
  exclude:
  - .git
  - web/sites/default/files
tasks:
  deploy:
    - command: drush updb
`
	if err := os.WriteFile(filepath.Join(root, ".ironstar", "config.yml"), []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)

	res, err := MigrateToIronignore()
	if err != nil {
		t.Fatal(err)
	}
	if !res.RemovedPackage {
		t.Fatal("expected RemovedPackage to be true")
	}
	if len(res.Ported) != 2 {
		t.Fatalf("expected 2 ported rules, got %v", res.Ported)
	}

	ig, err := os.ReadFile(filepath.Join(root, ".ironstar", ".ironignore"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(ig)
	for _, want := range []string{
		"Ported from config.yml package.exclude", // the dedicated group
		"web/sites/default/files",                // a ported rule
		"node_modules",                           // a default
	} {
		if !strings.Contains(s, want) {
			t.Fatalf(".ironignore missing %q:\n%s", want, s)
		}
	}

	newCfg, err := os.ReadFile(filepath.Join(root, ".ironstar", "config.yml"))
	if err != nil {
		t.Fatal(err)
	}
	cs := string(newCfg)
	if strings.Contains(cs, "package:") {
		t.Fatalf("config.yml still contains a package block:\n%s", cs)
	}
	for _, keep := range []string{"tasks:", "alias: demo", "drush updb"} {
		if !strings.Contains(cs, keep) {
			t.Fatalf("config.yml lost %q during migration:\n%s", keep, cs)
		}
	}

	// Re-running must refuse rather than clobber the existing .ironignore.
	if _, err := MigrateToIronignore(); err == nil {
		t.Fatal("expected an error when .ironignore already exists")
	}
}
