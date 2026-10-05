package motionsource

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveSetPathsRejectsPartialFamily(t *testing.T) {
	root := t.TempDir()
	set00Path := filepath.Join(root, "animal-source-set00.png")
	writeTestFile(t, set00Path)
	writeTestFile(t, filepath.Join(root, "animal-source-set01.png"))

	_, err := ResolveSetPaths(set00Path, 10)
	if err == nil {
		t.Fatal("ResolveSetPaths() accepted a partial set family")
	}
	for _, want := range []string{"incomplete motion source family", "found 2 of 10", "set02"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("ResolveSetPaths() error = %v, want %q", err, want)
		}
	}
}

func TestResolveSetPathsAllowsCanonicalSet00Only(t *testing.T) {
	root := t.TempDir()
	set00Path := filepath.Join(root, "animal-source-set00.png")
	writeTestFile(t, set00Path)

	paths, err := ResolveSetPaths(set00Path, 10)
	if err != nil {
		t.Fatalf("ResolveSetPaths() error = %v", err)
	}
	if len(paths) != 1 || paths[0] != set00Path {
		t.Fatalf("ResolveSetPaths() = %v, want only %q", paths, set00Path)
	}
}

func TestResolveSetPathsRequiresStandaloneSourceToExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "standalone.png")
	if _, err := ResolveSetPaths(path, 10); !os.IsNotExist(err) {
		t.Fatalf("ResolveSetPaths() error = %v, want missing-file error", err)
	}
}

func writeTestFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
