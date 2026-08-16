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

func TestBatches(t *testing.T) {
	tests := []struct {
		name  string
		names []string
		size  int
		want  [][]string
	}{
		{
			name:  "empty",
			names: nil,
			size:  20,
			want:  nil,
		},
		{
			name:  "fewer than size",
			names: []string{"a", "b", "c"},
			size:  20,
			want:  [][]string{{"a", "b", "c"}},
		},
		{
			name:  "exact multiple",
			names: []string{"a", "b", "c", "d"},
			size:  2,
			want:  [][]string{{"a", "b"}, {"c", "d"}},
		},
		{
			name:  "uneven remainder",
			names: []string{"a", "b", "c", "d", "e"},
			size:  2,
			want:  [][]string{{"a", "b"}, {"c", "d"}, {"e"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := batches(tt.names, tt.size)
			if len(got) != len(tt.want) {
				t.Fatalf("batches = %d groups, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if !testutil.EqualStrings(got[i], tt.want[i]) {
					t.Errorf("batches[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestAbsPaths(t *testing.T) {
	dir := "RAW"
	got, err := absPaths(dir, []string{"img1.cr2", "img2.cr2"})
	if err != nil {
		t.Fatalf("absPaths: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}

	want := []string{
		filepath.Join(cwd, "RAW", "img1.cr2"),
		filepath.Join(cwd, "RAW", "img2.cr2"),
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("absPaths[%d] = %q, want %q", i, got[i], want[i])
		}
		if !filepath.IsAbs(got[i]) {
			t.Errorf("absPaths[%d] = %q is not absolute", i, got[i])
		}
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
