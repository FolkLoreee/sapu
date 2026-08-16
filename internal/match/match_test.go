package match

import (
	"testing"

	"github.com/FolkLoreee/sapu/internal/testutil"
)

func TestClassify(t *testing.T) {
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
			keep, del := Classify(tt.jpegs, tt.raws)
			if !testutil.EqualStrings(keep, tt.wantKeep) {
				t.Errorf("keep = %v, want %v", keep, tt.wantKeep)
			}
			if !testutil.EqualStrings(del, tt.wantDel) {
				t.Errorf("del = %v, want %v", del, tt.wantDel)
			}
		})
	}
}
