package filetype

import "testing"

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
			got := Stem(tt.filename)
			if got != tt.want {
				t.Errorf("Stem(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}

func TestIsJPEG(t *testing.T) {
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
		if got := IsJPEG(tt.filename); got != tt.want {
			t.Errorf("IsJPEG(%q) = %v, want %v", tt.filename, got, tt.want)
		}
	}
}

func TestIsRaw(t *testing.T) {
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
		if got := IsRaw(tt.filename); got != tt.want {
			t.Errorf("IsRaw(%q) = %v, want %v", tt.filename, got, tt.want)
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
		if got := IsDotfile(tt.filename); got != tt.want {
			t.Errorf("IsDotfile(%q) = %v, want %v", tt.filename, got, tt.want)
		}
	}
}
