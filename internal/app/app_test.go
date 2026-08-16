package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/FolkLoreee/sapu/internal/testutil"
)

func TestRunHard(t *testing.T) {
	jpgDir := t.TempDir()
	rawDir := t.TempDir()

	testutil.WriteFile(t, filepath.Join(jpgDir, "img1.jpg"))
	testutil.WriteFile(t, filepath.Join(jpgDir, "img5.jpg"))

	testutil.WriteFile(t, filepath.Join(rawDir, "img1.cr2"))
	testutil.WriteFile(t, filepath.Join(rawDir, "img2.cr2"))
	testutil.WriteFile(t, filepath.Join(rawDir, "img5.cr2"))

	if err := Run(jpgDir, rawDir, true, false, os.Stdout); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if _, err := os.Stat(filepath.Join(rawDir, "img1.cr2")); err != nil {
		t.Errorf("img1.cr2 should be kept: %v", err)
	}
	if _, err := os.Stat(filepath.Join(rawDir, "img5.cr2")); err != nil {
		t.Errorf("img5.cr2 should be kept: %v", err)
	}
	if _, err := os.Stat(filepath.Join(rawDir, "img2.cr2")); !os.IsNotExist(err) {
		t.Errorf("img2.cr2 should be removed: %v", err)
	}
}

func TestRunSameDir(t *testing.T) {
	dir := t.TempDir()
	if err := Run(dir, dir, true, false, os.Stdout); err == nil {
		t.Fatal("expected error for same dir, got nil")
	}
}

func TestRunMissingDir(t *testing.T) {
	if err := Run("/nonexistent", t.TempDir(), true, false, os.Stdout); err == nil {
		t.Fatal("expected error for missing dir, got nil")
	}
}

func TestRunDryRun(t *testing.T) {
	jpgDir := t.TempDir()
	rawDir := t.TempDir()

	testutil.WriteFile(t, filepath.Join(jpgDir, "img1.jpg"))
	testutil.WriteFile(t, filepath.Join(rawDir, "img1.cr2"))
	testutil.WriteFile(t, filepath.Join(rawDir, "img2.cr2"))

	var buf bytes.Buffer
	if err := Run(jpgDir, rawDir, false, true, &buf); err != nil {
		t.Fatalf("Run: %v", err)
	}

	for _, name := range []string{"img1.cr2", "img2.cr2"} {
		if _, err := os.Stat(filepath.Join(rawDir, name)); err != nil {
			t.Errorf("%s should remain after dry-run: %v", name, err)
		}
	}
	if !bytes.Contains(buf.Bytes(), []byte("[DRY-RUN] would remove:")) {
		t.Errorf("dry-run output missing preview line:\n%s", buf.String())
	}
}
