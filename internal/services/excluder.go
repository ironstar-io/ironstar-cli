package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/ironstar-cli/cmd/flags"
	"github.com/ironstar-io/ironstar-cli/internal/system/tarball"
	"github.com/ironstar-io/ironstar-cli/internal/types"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

// gitIgnoreExcluder adapts a go-git gitignore matcher to tarball.Excluder so
// .ironstar/.ironignore is honoured with full git semantics (**, negation,
// anchoring, directory-only).
type gitIgnoreExcluder struct {
	matcher gitignore.Matcher
}

func (g gitIgnoreExcluder) Excludes(relPath string, isDir bool) (bool, error) {
	if relPath == "" || relPath == "." {
		return false, nil
	}
	return g.matcher.Match(strings.Split(relPath, "/"), isDir), nil
}

func newGitIgnoreExcluder(ignoreFilePath string, extraPatterns []string) (tarball.Excluder, error) {
	b, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		return nil, err
	}

	var patterns []gitignore.Pattern
	for line := range strings.SplitSeq(string(b), "\n") {
		s := strings.TrimSuffix(line, "\r")
		if strings.HasPrefix(s, "#") || strings.TrimSpace(s) == "" {
			continue
		}
		patterns = append(patterns, gitignore.ParsePattern(s, nil))
	}

	// --exclude patterns are appended last so they take precedence.
	for _, p := range extraPatterns {
		patterns = append(patterns, gitignore.ParsePattern(p, nil))
	}

	return gitIgnoreExcluder{matcher: gitignore.NewMatcher(patterns)}, nil
}

// excludeSource names which exclude ruleset packaging applied, so the choice can
// be surfaced to the user. label is empty when nothing excludes files.
type excludeSource struct {
	label string
	both  bool // .ironignore is in effect while config.yml package.exclude is also set
}

// resolvePackageExcluder picks the exclude source for packaging:
// .ironstar/.ironignore (full gitignore semantics) when present, otherwise the
// config.yml package.exclude list (root-anchored slashes, basename-at-any-depth
// for slash-less patterns). It is the one place both the real tarball and the
// dry-run index resolve excludes; the returned excludeSource lets callers report
// which ruleset was used.
func resolvePackageExcluder(projectRoot string, proj types.ProjectConfig, flg flags.Accumulator) (tarball.Excluder, excludeSource, error) {
	extra := splitCSV(flg.Exclude)
	ignorePath := filepath.Join(projectRoot, ".ironstar", ".ironignore")

	info, err := os.Stat(ignorePath)
	switch {
	case err == nil && !info.IsDir():
		ex, gerr := newGitIgnoreExcluder(ignorePath, extra)
		return ex, excludeSource{label: ".ironstar/.ironignore", both: len(proj.Package.Exclude) > 0}, gerr
	case err != nil && !os.IsNotExist(err):
		return nil, excludeSource{}, err
	}

	patterns := append(append([]string{}, proj.Package.Exclude...), extra...)
	var src excludeSource
	switch {
	case len(proj.Package.Exclude) > 0:
		src.label = ".ironstar/config.yml (package.exclude)"
	case len(extra) > 0:
		src.label = "--exclude flag"
	}
	return tarball.NewPatternExcluder(patterns), src, nil
}

// printUploadNotice tells the user, in neutral terms, that packaging uploads the
// project directory and which exclude ruleset was applied.
func printUploadNotice(src excludeSource) {
	if src.label != "" {
		fmt.Printf("This will upload your project directory to Ironstar, excluding paths in %s.\n", src.label)
	} else {
		fmt.Println("This will upload your project directory to Ironstar. No exclude rules are configured.")
	}
	fmt.Println("Files that aren't excluded (e.g. .env files or database dumps) will be included — review before continuing.")
	if src.both {
		color.Yellow("Note: .ironstar/.ironignore takes precedence; the package.exclude list in config.yml is ignored.")
	}
}

func splitCSV(s string) []string {
	out := []string{}
	for p := range strings.SplitSeq(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
