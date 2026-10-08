package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStaticFilePathAllowsNormalFiles(t *testing.T) {
	base := t.TempDir()
	got, ok := staticFilePath(base, "/assets/app.js")
	if !ok {
		t.Fatal("expected ok for normal path")
	}
	absBase, _ := filepath.Abs(base)
	if want := filepath.Join(absBase, "assets", "app.js"); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStaticFilePathRejectsDotDotNames(t *testing.T) {
	base := t.TempDir()
	if _, ok := staticFilePath(base, "/foo..bar"); ok {
		t.Error("expected rejection for name containing '..'")
	}
}

func TestStaticFilePathNeverEscapesBase(t *testing.T) {
	base := t.TempDir()
	absBase, _ := filepath.Abs(base)
	attacks := []string{
		"/../../etc/passwd",
		"/..",
		"/../..",
		"/assets/../../../etc/shadow",
		"/..%2F..%2Fetc/passwd",
		"/..\\..\\windows\\system32",
	}
	for _, p := range attacks {
		got, ok := staticFilePath(base, p)
		if !ok {
			continue
		}
		if got != absBase && !strings.HasPrefix(got, absBase+string(os.PathSeparator)) {
			t.Errorf("staticFilePath(%q) = %q escapes base %q", p, got, absBase)
		}
	}
}
