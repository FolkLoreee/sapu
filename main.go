package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var jpegExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
}

var rawExts = map[string]bool{
	".cr2": true,
	".cr3": true,
	".nef": true,
	".nrw": true,
	".arw": true,
	".dng": true,
	".raf": true,
	".orf": true,
	".rw2": true,
	".pef": true,
	".srw": true,
}

func stem(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		return filename
	}
	return filename[:idx]
}

func isJpegExt(filename string) bool {
	return jpegExts[strings.ToLower(filepath.Ext(filename))]
}

func isRawExt(filename string) bool {
	return rawExts[strings.ToLower(filepath.Ext(filename))]
}

func isDotfile(filename string) bool {
	return strings.HasPrefix(filename, ".")
}

func filesToRemove(jpegs, raws []string) (keep, del []string) {
	jpegStems := make(map[string]bool, len(jpegs))
	for _, f := range jpegs {
		if isDotfile(f) || !isJpegExt(f) {
			continue
		}
		jpegStems[stem(f)] = true
	}

	for _, f := range raws {
		if isDotfile(f) || !isRawExt(f) {
			continue
		}
		if jpegStems[stem(f)] {
			keep = append(keep, f)
		} else {
			del = append(del, f)
		}
	}
	return keep, del
}

func listFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if isDotfile(e.Name()) {
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

func removeFiles(dir string, names []string, hard, dryRun bool) error {
	if len(names) == 0 {
		return nil
	}
	if dryRun {
		for _, name := range names {
			fmt.Printf("[DRY-RUN] would remove: %s\n", filepath.Join(dir, name))
		}
		return nil
	}
	if hard {
		for _, name := range names {
			path := filepath.Join(dir, name)
			fmt.Printf("remove: %s\n", path)
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	}
	return trashFiles(dir, names)
}

func trashFiles(dir string, names []string) error {
	paths := make([]string, len(names))
	for i, name := range names {
		paths[i] = filepath.Join(dir, name)
		fmt.Printf("remove: %s\n", paths[i])
	}
	script := trashScript(paths)
	cmd := exec.Command("osascript", "-e", script)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func trashScript(paths []string) string {
	items := make([]string, len(paths))
	for i, p := range paths {
		items[i] = fmt.Sprintf(`POSIX file "%s"`, strings.ReplaceAll(p, `"`, `\"`))
	}
	return fmt.Sprintf(`tell application "Finder" to delete {%s}`, strings.Join(items, ", "))
}

func run(jpgDir, rawDir string, hard, dryRun bool) error {
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

	jpegs, err := listFiles(jpgDir)
	if err != nil {
		return fmt.Errorf("list jpeg dir: %w", err)
	}
	raws, err := listFiles(rawDir)
	if err != nil {
		return fmt.Errorf("list raw dir: %w", err)
	}

	keep, del := filesToRemove(jpegs, raws)

	if err := removeFiles(rawDir, del, hard, dryRun); err != nil {
		return err
	}

	fmt.Printf("scanned %d raw files: %d kept, %d removed\n", len(keep)+len(del), len(keep), len(del))
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

func main() {
	jpgDir := flag.String("jpg", "", "directory containing jpeg files")
	rawDir := flag.String("raw", "", "directory containing raw files")
	hard := flag.Bool("hard", false, "permanently delete unmatched raw files instead of moving to Trash")
	dryRun := flag.Bool("dry-run", false, "show what would be removed without doing anything")
	flag.Parse()

	if *jpgDir == "" || *rawDir == "" {
		fmt.Fprintln(os.Stderr, "both --jpg and --raw are required")
		flag.Usage()
		os.Exit(2)
	}

	if err := run(*jpgDir, *rawDir, *hard, *dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "sapu: %v\n", err)
		os.Exit(1)
	}
}
