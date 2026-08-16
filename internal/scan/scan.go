// Package scan lists the files in a directory.
package scan

import (
	"os"

	"github.com/FolkLoreee/sapu/internal/filetype"
)

// List returns the names of all regular, non-dotfile entries in dir, sorted
// lexically. Subdirectories and symlinks are skipped.
func List(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filetype.IsDotfile(e.Name()) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			return nil, err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			continue
		}
		files = append(files, e.Name())
	}
	return files, nil
}
