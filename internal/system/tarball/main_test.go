package tarball

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

// writeTree lays down a small project tree and returns the root.
func writeTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()

	files := map[string]string{
		"a.txt":                         "abc",       // 3   kept
		"sub/b.txt":                     "hello",     // 5   kept
		"vendor/lib.txt":                "vendored",  // 8   kept
		"web/core/index.php":            "php11",     // 5   kept
		".git/config":                   "gitstuff",  // excluded (root .git)
		"web/core/.git/packed":          "pack",      // excluded (NESTED .git, slash-less pattern)
		"web/sites/default/files/m.bin": "mediadata", // excluded (anchored slashed pattern)
	}
	for rel, body := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func paths(entries []IndexEntry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Path
	}
	sort.Strings(out)
	return out
}

func TestIndexHonorsExcludes(t *testing.T) {
	root := writeTree(t)
	excludes := []string{".git", "web/sites/default/files"}

	entries, total, err := Index(root, NewPatternExcluder(excludes))
	if err != nil {
		t.Fatal(err)
	}

	sizes := map[string]int64{}
	for _, e := range entries {
		sizes[e.Path] = e.Size
	}

	// Excluded entries must be absent (a real check: the subjects exist on disk).
	// web/core/.git is excluded by the slash-less `.git` matching at depth.
	for _, gone := range []string{".git/config", "web/core/.git/packed", "web/sites/default/files/m.bin"} {
		if _, ok := sizes[gone]; ok {
			t.Fatalf("expected %q to be excluded, but it was indexed", gone)
		}
	}

	// Included entries must be present with correct uncompressed sizes.
	// web/core/index.php proves only the nested .git is dropped, not all of web/core.
	want := map[string]int64{"a.txt": 3, "sub/b.txt": 5, "vendor/lib.txt": 8, "web/core/index.php": 5}
	for p, sz := range want {
		got, ok := sizes[p]
		if !ok {
			t.Fatalf("expected %q to be indexed", p)
		}
		if got != sz {
			t.Fatalf("size for %q: want %d, got %d", p, sz, got)
		}
	}

	if len(entries) != len(want) {
		t.Fatalf("expected %d files, got %d (%v)", len(want), len(entries), paths(entries))
	}
	if total != 21 {
		t.Fatalf("expected total 21 bytes, got %d", total)
	}
}

func TestPatternExcluderAnchoringAndDepth(t *testing.T) {
	ex := NewPatternExcluder([]string{".git", "node_modules", "*.sql", "web/sites/default/files"})

	cases := []struct {
		path  string
		isDir bool
		want  bool
	}{
		{".git", true, true},                          // slash-less, root
		{"web/core/.git", true, true},                 // slash-less, nested
		{"web/themes/x/node_modules", true, true},     // slash-less, deeply nested
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

// TestIndexMatchesNewTarGZ guards the shared Excluder: the dry-run index must
// list exactly the regular files NewTarGZ packages.
func TestIndexMatchesNewTarGZ(t *testing.T) {
	root := writeTree(t)
	excludes := []string{".git", "web/sites/default/files"}

	tarPath := filepath.Join(t.TempDir(), "out.tar.gz")
	if err := NewTarGZ(tarPath, root, NewPatternExcluder(excludes)); err != nil {
		t.Fatal(err)
	}

	archived, _, err := IndexArchive(tarPath)
	if err != nil {
		t.Fatal(err)
	}

	indexed, _, err := Index(root, NewPatternExcluder(excludes))
	if err != nil {
		t.Fatal(err)
	}

	gotArchive := paths(archived)
	gotIndex := paths(indexed)
	if len(gotArchive) != len(gotIndex) {
		t.Fatalf("packaged set %v != indexed set %v", gotArchive, gotIndex)
	}
	for i := range gotArchive {
		if gotArchive[i] != gotIndex[i] {
			t.Fatalf("packaged set %v != indexed set %v", gotArchive, gotIndex)
		}
	}
}

func TestIndexArchiveRoundTrip(t *testing.T) {
	root := writeTree(t)
	excludes := []string{".git", "web/sites/default/files"}

	tarPath := filepath.Join(t.TempDir(), "out.tar.gz")
	if err := NewTarGZ(tarPath, root, NewPatternExcluder(excludes)); err != nil {
		t.Fatal(err)
	}

	entries, total, err := IndexArchive(tarPath)
	if err != nil {
		t.Fatal(err)
	}

	// Cross-check against a direct read of the archive's regular files.
	wantSizes := map[string]int64{}
	var wantTotal int64
	f, err := os.Open(tarPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if hdr.Typeflag == tar.TypeReg {
			wantSizes[hdr.Name] = hdr.Size
			wantTotal += hdr.Size
		}
	}

	if total != wantTotal {
		t.Fatalf("total: want %d, got %d", wantTotal, total)
	}
	if len(entries) != len(wantSizes) {
		t.Fatalf("entry count: want %d, got %d", len(wantSizes), len(entries))
	}
	for _, e := range entries {
		if wantSizes[e.Path] != e.Size {
			t.Fatalf("size for %q: want %d, got %d", e.Path, wantSizes[e.Path], e.Size)
		}
	}
}
