// Package filetype classifies image files by extension and derives their
// stems for matching.
package filetype

import (
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

// Stem returns filename with its final extension removed. A filename without
// a dot is returned unchanged.
func Stem(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		return filename
	}
	return filename[:idx]
}

// IsJPEG reports whether filename has a .jpg or .jpeg extension,
// case-insensitively.
func IsJPEG(filename string) bool {
	return jpegExts[strings.ToLower(filepath.Ext(filename))]
}

// IsRaw reports whether filename has a supported raw extension,
// case-insensitively.
func IsRaw(filename string) bool {
	return rawExts[strings.ToLower(filepath.Ext(filename))]
}

// IsDotfile reports whether filename begins with a dot.
func IsDotfile(filename string) bool {
	return strings.HasPrefix(filename, ".")
}
