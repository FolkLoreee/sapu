// Package app wires scanning, matching, and removal together into the sapu
// command's behavior.
package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/FolkLoreee/sapu/internal/match"
	"github.com/FolkLoreee/sapu/internal/remove"
	"github.com/FolkLoreee/sapu/internal/scan"
)

// Run scans jpgDir for JPEG stems and removes every raw file in rawDir whose
// stem does not match. When dryRun is true it only reports what would be
// removed. All output is written to out.
func Run(jpgDir, rawDir string, hard, dryRun bool, out io.Writer) error {
	jpgInfo, err := os.Stat(jpgDir)
	if err != nil {
		return fmt.Errorf("jpeg dir: %w", err)
	}
	if !jpgInfo.IsDir() {
		return fmt.Errorf("jpeg dir %q is not a directory", jpgDir)
	}

	rawInfo, err := os.Stat(rawDir)
	if err != nil {
		return fmt.Errorf("raw dir: %w", err)
	}
	if !rawInfo.IsDir() {
		return fmt.Errorf("raw dir %q is not a directory", rawDir)
	}

	if samePath(jpgDir, rawDir) {
		return fmt.Errorf("jpeg dir and raw dir must be different")
	}

	jpegs, err := scan.List(jpgDir)
	if err != nil {
		return fmt.Errorf("list jpeg dir: %w", err)
	}
	raws, err := scan.List(rawDir)
	if err != nil {
		return fmt.Errorf("list raw dir: %w", err)
	}

	keep, del := match.Classify(jpegs, raws)

	if dryRun {
		for _, name := range del {
			fmt.Fprintf(out, "[DRY-RUN] would remove: %s\n", filepath.Join(rawDir, name))
		}
	} else {
		for _, name := range del {
			fmt.Fprintf(out, "remove: %s\n", filepath.Join(rawDir, name))
		}
		if err := remove.Files(rawDir, del, hard); err != nil {
			return err
		}
	}

	fmt.Fprintf(out, "scanned %d raw files: %d kept, %d removed\n", len(keep)+len(del), len(keep), len(del))
	return nil
}

func samePath(a, b string) bool {
	aa, errA := filepath.Abs(a)
	bb, errB := filepath.Abs(b)
	if errA != nil || errB != nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	return aa == bb
}
