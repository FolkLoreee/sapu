// Package match classifies raw files as keep or delete based on whether their
// stem matches a JPEG stem.
package match

import "github.com/FolkLoreee/sapu/internal/filetype"

// Classify splits raws into files to keep (stem matches a JPEG) and files to
// delete (no matching JPEG). Dotfiles and files without recognized JPEG/raw
// extensions are ignored.
func Classify(jpegs, raws []string) (keep, del []string) {
	jpegStems := make(map[string]bool, len(jpegs))
	for _, f := range jpegs {
		if filetype.IsDotfile(f) || !filetype.IsJPEG(f) {
			continue
		}
		jpegStems[filetype.Stem(f)] = true
	}

	for _, f := range raws {
		if filetype.IsDotfile(f) || !filetype.IsRaw(f) {
			continue
		}
		if jpegStems[filetype.Stem(f)] {
			keep = append(keep, f)
		} else {
			del = append(del, f)
		}
	}
	return keep, del
}
