package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStem(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{name: "simple jpg", filename: "img1.jpg", want: "img1"},
		{name: "simple raw", filename: "img1.cr2", want: "img1"},
		{name: "multiple dots", filename: "img.v1.jpg", want: "img.v1"},
		{name: "no extension", filename: "README", want: "README"},
		{name: "hidden file with ext", filename: ".hidden.jpg", want: ".hidden"},
		{name: "uppercase ext", filename: "IMG.CR2", want: "IMG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stem(tt.filename)
			if got != tt.want {
				t.Errorf("stem(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}

func TestIsJpegExt(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		{"img1.jpg", true},
		{"img1.JPG", true},
		{"img1.jpeg", true},
		{"img1.Jpeg", true},
		{"img1.png", false},
		{"img1.cr2", false},
		{"img1", false},
	}

	for _, tt := range tests {
		if got := isJpegExt(tt.filename); got != tt.want {
			t.Errorf("isJpegExt(%q) = %v, want %v", tt.filename, got, tt.want)
		}
	}
}

func TestIsRawExt(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		{"img1.cr2", true},
		{"img1.CR2", true},
		{"img1.cr3", true},
		{"img1.nef", true},
		{"img1.nrw", true},
		{"img1.arw", true},
		{"img1.dng", true},
		{"img1.raf", true},
		{"img1.orf", true},
		{"img1.rw2", true},
		{"img1.pef", true},
		{"img1.srw", true},
		{"img1.jpg", false},
		{"img1.png", false},
		{"img1", false},
	}

	for _, tt := range tests {
		if got := isRawExt(tt.filename); got != tt.want {
			t.Errorf("isRawExt(%q) = %v, want %v", tt.filename, got, tt.want)
		}
	}
}

func TestIsDotfile(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		{".DS_Store", true},
		{".hidden", true},
		{"img1.jpg", false},
		{"..", true},
		{".", true},
	}

	for _, tt := range tests {
		if got := isDotfile(tt.filename); got != tt.want {
			t.Errorf("isDotfile(%q) = %v, want %v", tt.filename, got, tt.want)
		}
	}
}

func TestFilesToRemove(t *testing.T) {
	tests := []struct {
		name     string
		jpegs    []string
		raws     []string
		wantKeep []string
		wantDel  []string
	}{
		{
			name:     "example from spec",
			jpegs:    []string{"img1.jpg", "img5.jpg"},
			raws:     []string{"img1.cr2", "img2.cr2", "img3.cr2", "img4.cr2", "img5.cr2"},
			wantKeep: []string{"img1.cr2", "img5.cr2"},
			wantDel:  []string{"img2.cr2", "img3.cr2", "img4.cr2"},
		},
		{
			name:     "case-sensitive stems",
			jpegs:    []string{"img1.jpg"},
			raws:     []string{"IMG1.cr2"},
			wantKeep: []string{},
			wantDel:  []string{"IMG1.cr2"},
		},
		{
			name:     "case-insensitive extensions",
			jpegs:    []string{"IMG1.JPG"},
			raws:     []string{"IMG1.cr2"},
			wantKeep: []string{"IMG1.cr2"},
			wantDel:  []string{},
		},
		{
			name:     "multiple raw files per stem kept",
			jpegs:    []string{"img1.jpg"},
			raws:     []string{"img1.cr2", "img1.dng"},
			wantKeep: []string{"img1.cr2", "img1.dng"},
			wantDel:  []string{},
		},
		{
			name:     "multi-dot stems",
			jpegs:    []string{"img.v1.jpg"},
			raws:     []string{"img.v1.cr2", "img.cr2"},
			wantKeep: []string{"img.v1.cr2"},
			wantDel:  []string{"img.cr2"},
		},
		{
			name:     "non-raw extensions untouched",
			jpegs:    []string{"img1.jpg"},
			raws:     []string{"img1.cr2", "img2.txt", "img2.cr2"},
			wantKeep: []string{"img1.cr2"},
			wantDel:  []string{"img2.cr2"},
		},
		{
			name:     "dotfiles skipped",
			jpegs:    []string{"img1.jpg"},
			raws:     []string{"img1.cr2", ".DS_Store", "img2.cr2"},
			wantKeep: []string{"img1.cr2"},
			wantDel:  []string{"img2.cr2"},
		},
		{
			name:     "no raws",
			jpegs:    []string{"img1.jpg"},
			raws:     []string{},
			wantKeep: []string{},
			wantDel:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keep, del := filesToRemove(tt.jpegs, tt.raws)
			if !equalStrings(keep, tt.wantKeep) {
				t.Errorf("keep = %v, want %v", keep, tt.wantKeep)
			}
			if !equalStrings(del, tt.wantDel) {
				t.Errorf("del = %v, want %v", del, tt.wantDel)
			}
		})
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestListFiles(t *testing.T) {
	dir := t.TempDir()

	writeTestFile(t, filepath.Join(dir, "img1.jpg"))
	writeTestFile(t, filepath.Join(dir, "img2.cr2"))
	writeTestFile(t, filepath.Join(dir, ".DS_Store"))
	os.MkdirAll(filepath.Join(dir, "subdir"), 0o755)
	writeTestFile(t, filepath.Join(dir, "subdir", "nested.jpg"))

	got, err := listFiles(dir)
	if err != nil {
		t.Fatalf("listFiles: %v", err)
	}

	want := []string{"img1.jpg", "img2.cr2"}
	if !equalStrings(got, want) {
		t.Errorf("listFiles = %v, want %v", got, want)
	}
}

func TestListFilesMissingDir(t *testing.T) {
	_, err := listFiles("/nonexistent/dir/path")
	if err == nil {
		t.Fatal("expected error for missing dir, got nil")
	}
}

func TestListFilesSymlinkIgnored(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "real.cr2")
	writeTestFile(t, target)
	os.Symlink(target, filepath.Join(dir, "link.cr2"))

	got, err := listFiles(dir)
	if err != nil {
		t.Fatalf("listFiles: %v", err)
	}
	want := []string{"real.cr2"}
	if !equalStrings(got, want) {
		t.Errorf("listFiles = %v, want %v", got, want)
	}
}

func writeTestFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}

func TestRemoveFilesHard(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "img1.cr2")
	f2 := filepath.Join(dir, "img2.cr2")
	writeTestFile(t, f1)
	writeTestFile(t, f2)

	if err := removeFiles(dir, []string{"img1.cr2", "img2.cr2"}, true, false); err != nil {
		t.Fatalf("removeFiles: %v", err)
	}

	for _, f := range []string{f1, f2} {
		if _, err := os.Stat(f); !os.IsNotExist(err) {
			t.Errorf("expected %s to be removed, stat err = %v", f, err)
		}
	}
}

func TestRemoveFilesDryRun(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "img1.cr2")
	writeTestFile(t, f1)

	if err := removeFiles(dir, []string{"img1.cr2"}, true, true); err != nil {
		t.Fatalf("removeFiles: %v", err)
	}

	if _, err := os.Stat(f1); err != nil {
		t.Errorf("dry-run should not remove file, stat err = %v", err)
	}
}

func TestTrashScript(t *testing.T) {
	paths := []string{"/photos/raw/img2.cr2", "/photos/raw/img3.cr2"}
	script := trashScript(paths)

	want := `tell application "Finder" to delete {POSIX file "/photos/raw/img2.cr2", POSIX file "/photos/raw/img3.cr2"}`
	if script != want {
		t.Errorf("trashScript = %q, want %q", script, want)
	}
}

func TestRunHard(t *testing.T) {
	jpgDir := t.TempDir()
	rawDir := t.TempDir()

	writeTestFile(t, filepath.Join(jpgDir, "img1.jpg"))
	writeTestFile(t, filepath.Join(jpgDir, "img5.jpg"))

	writeTestFile(t, filepath.Join(rawDir, "img1.cr2"))
	writeTestFile(t, filepath.Join(rawDir, "img2.cr2"))
	writeTestFile(t, filepath.Join(rawDir, "img5.cr2"))

	if err := run(jpgDir, rawDir, true, false); err != nil {
		t.Fatalf("run: %v", err)
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
	if err := run(dir, dir, true, false); err == nil {
		t.Fatal("expected error for same dir, got nil")
	}
}

func TestRunMissingDir(t *testing.T) {
	if err := run("/nonexistent", t.TempDir(), true, false); err == nil {
		t.Fatal("expected error for missing dir, got nil")
	}
}
