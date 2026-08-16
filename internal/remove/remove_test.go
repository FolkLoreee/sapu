package remove

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/FolkLoreee/sapu/internal/testutil"
)

func TestFilesHard(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "img1.cr2")
	f2 := filepath.Join(dir, "img2.cr2")
	testutil.WriteFile(t, f1)
	testutil.WriteFile(t, f2)

	if err := Files(dir, []string{"img1.cr2", "img2.cr2"}, true); err != nil {
		t.Fatalf("Files: %v", err)
	}

	for _, f := range []string{f1, f2} {
		if _, err := os.Stat(f); !os.IsNotExist(err) {
			t.Errorf("expected %s to be removed, stat err = %v", f, err)
		}
	}
}

func TestFilesEmpty(t *testing.T) {
	dir := t.TempDir()
	if err := Files(dir, nil, true); err != nil {
		t.Fatalf("Files with no names: %v", err)
	}
}

func TestTrashScript(t *testing.T) {
	paths := []string{"/photos/raw/img2.cr2", "/photos/raw/img3.cr2"}
	script := TrashScript(paths)

	want := `tell application "Finder" to delete {POSIX file "/photos/raw/img2.cr2", POSIX file "/photos/raw/img3.cr2"}`
	if script != want {
		t.Errorf("TrashScript = %q, want %q", script, want)
	}
}
