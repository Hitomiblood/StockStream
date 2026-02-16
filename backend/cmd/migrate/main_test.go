package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestFileSourceURL(t *testing.T) {
	url := fileSourceURL("C:/tmp/migrations")
	if !strings.HasPrefix(url, "file://") {
		t.Fatalf("unexpected file source url: %s", url)
	}
	if !strings.Contains(url, "/migrations") {
		t.Fatalf("url should include path, got: %s", url)
	}
}

func TestFindMigrationsDir_PrefersLocalFolder(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}

	tmp := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmp, "migrations"), 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}

	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir tmp: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	got := findMigrationsDir()
	if !strings.Contains(filepath.ToSlash(got), "/migrations") {
		t.Fatalf("findMigrationsDir unexpected: %s", got)
	}
}

func TestFindMigrationsDir_FallbackBackendFolder(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}

	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, "backend", "migrations"), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("Chdir tmp: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	got := findMigrationsDir()
	if runtime.GOOS == "windows" {
		if !strings.Contains(strings.ToLower(got), "backend") {
			t.Fatalf("expected backend fallback path, got: %s", got)
		}
	} else if !strings.Contains(got, "/backend/migrations") {
		t.Fatalf("expected backend fallback path, got: %s", got)
	}
}
