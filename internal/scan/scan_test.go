package scan

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/FolkLoreee/sapu/internal/testutil"
)

func TestList(t *testing.T) {
	dir := t.TempDir()

	testutil.WriteFile(t, filepath.Join(dir, "img1.jpg"))
	testutil.WriteFile(t, filepath.Join(dir, "img2.cr2"))
	testutil.WriteFile(t, filepath.Join(dir, ".DS_Store"))
	os.MkdirAll(filepath.Join(dir, "subdir"), 0o755)
	testutil.WriteFile(t, filepath.Join(dir, "subdir", "nested.jpg"))

	got, err := List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	want := []string{"img1.jpg", "img2.cr2"}
	if !testutil.EqualStrings(got, want) {
		t.Errorf("List = %v, want %v", got, want)
	}
}

func TestListMissingDir(t *testing.T) {
	_, err := List("/nonexistent/dir/path")
	if err == nil {
		t.Fatal("expected error for missing dir, got nil")
	}
}

func TestListSymlinkIgnored(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "real.cr2")
	testutil.WriteFile(t, target)
	os.Symlink(target, filepath.Join(dir, "link.cr2"))

	got, err := List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	want := []string{"real.cr2"}
	if !testutil.EqualStrings(got, want) {
		t.Errorf("List = %v, want %v", got, want)
	}
}
