// Package remove deletes raw files from disk, either permanently or to the
// macOS Trash.
package remove

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Files deletes the named files from dir. When hard is true files are removed
// permanently; otherwise they are moved to the macOS Trash via Finder.
func Files(dir string, names []string, hard bool) error {
	if len(names) == 0 {
		return nil
	}
	if hard {
		for _, name := range names {
			if err := os.Remove(filepath.Join(dir, name)); err != nil {
				return err
			}
		}
		return nil
	}

	for _, batch := range batches(names, trashBatchSize) {
		paths, err := absPaths(dir, batch)
		if err != nil {
			return err
		}
		cmd := exec.Command("osascript", "-e", TrashScript(paths))
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

// absPaths returns the absolute paths of the named files under dir. Finder's
// delete command cannot resolve relative POSIX paths.
func absPaths(dir string, names []string) ([]string, error) {
	paths := make([]string, len(names))
	for i, name := range names {
		abs, err := filepath.Abs(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		paths[i] = abs
	}
	return paths, nil
}

// trashBatchSize caps the number of files passed to Finder in a single delete
// command. Large lists make Finder intermittently fail with error -10010.
const trashBatchSize = 100

// batches splits names into consecutive groups of at most size elements,
// preserving order.
func batches(names []string, size int) [][]string {
	var out [][]string
	for i := 0; i < len(names); i += size {
		end := i + size
		if end > len(names) {
			end = len(names)
		}
		out = append(out, names[i:end])
	}
	return out
}

// TrashScript returns an AppleScript snippet that moves the given paths to the
// macOS Trash in a single Finder delete command.
func TrashScript(paths []string) string {
	items := make([]string, len(paths))
	for i, p := range paths {
		items[i] = fmt.Sprintf(`POSIX file "%s"`, strings.ReplaceAll(p, `"`, `\"`))
	}
	return fmt.Sprintf(`tell application "Finder" to delete {%s}`, strings.Join(items, ", "))
}
