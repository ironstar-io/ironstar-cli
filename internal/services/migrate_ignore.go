package services

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/ironstar-cli/internal/constants"
	"github.com/ironstar-io/ironstar-cli/internal/system/fs"

	"github.com/pkg/errors"
)

// IronignoreMigration summarises the result of `iron init ignore` so the command
// can report exactly what changed.
type IronignoreMigration struct {
	IgnorePath     string
	ConfigPath     string
	Ported         []string
	RemovedPackage bool
}

// MigrateToIronignore creates .ironstar/.ironignore from the existing
// config.yml package.exclude (ported into their own group) plus the standard
// defaults, then strips the package block from config.yml.
func MigrateToIronignore() (IronignoreMigration, error) {
	var res IronignoreMigration

	root := fs.ProjectRoot()
	if root == constants.ProjectRootNotFound {
		return res, errors.New("Could not locate the project root. Run this from inside your project.")
	}

	confPath := filepath.Join(root, ".ironstar", "config.yml")
	if !fs.CheckExists(confPath) {
		return res, errors.New("No .ironstar/config.yml found. Run `iron init` first.")
	}

	ignorePath := filepath.Join(root, ".ironstar", ".ironignore")
	if fs.CheckExists(ignorePath) {
		return res, errors.New("A .ironstar/.ironignore already exists. Remove it first if you want to regenerate it.")
	}

	proj, err := ReadInProjectConfig(root)
	if err != nil {
		return res, err
	}
	res.IgnorePath = ignorePath
	res.ConfigPath = confPath
	res.Ported = proj.Package.Exclude

	content := ironignoreHeader + "\n" + portedGroup(res.Ported) + "\n" + defaultIronignoreRules
	if err := fs.TouchByteArray(ignorePath, []byte(content), 0644); err != nil {
		return res, err
	}

	raw, err := os.ReadFile(confPath)
	if err != nil {
		return res, err
	}
	if newRaw, removed := removePackageBlock(raw); removed {
		fs.Replace(confPath, newRaw, 0400)
		res.RemovedPackage = true
	}

	return res, nil
}

// portedGroup renders the config.yml excludes as their own clearly labelled
// block so it is obvious which rules were migrated.
func portedGroup(excludes []string) string {
	var b strings.Builder
	b.WriteString("# --- Ported from config.yml package.exclude (by `iron init ignore`) ---\n")
	if len(excludes) == 0 {
		b.WriteString("# (no package.exclude was present in config.yml)\n")
		return b.String()
	}
	for _, e := range excludes {
		if e = strings.TrimSpace(e); e != "" {
			b.WriteString(e)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// removePackageBlock strips a top-level `package:` mapping (its key line plus the
// indented lines and blank lines that follow it) from YAML, leaving every other
// key, comment, and the original formatting untouched. It returns the original
// bytes unchanged when no package block is present.
func removePackageBlock(raw []byte) ([]byte, bool) {
	lines := strings.Split(string(raw), "\n")
	out := make([]string, 0, len(lines))
	removed := false

	i := 0
	for i < len(lines) {
		if isTopLevelKey(lines[i], "package") {
			removed = true
			i++
			for i < len(lines) && (lines[i] == "" || lines[i][0] == ' ' || lines[i][0] == '\t') {
				i++
			}
			continue
		}
		out = append(out, lines[i])
		i++
	}

	if !removed {
		return raw, false
	}
	return []byte(strings.Join(out, "\n")), true
}

// isTopLevelKey reports whether line declares `key:` at column 0 (a top-level
// mapping key), not a nested or commented occurrence.
func isTopLevelKey(line, key string) bool {
	if !strings.HasPrefix(line, key+":") {
		return false
	}
	rest := line[len(key)+1:]
	return rest == "" || rest[0] == ' ' || rest[0] == '\t' || rest[0] == '{'
}
