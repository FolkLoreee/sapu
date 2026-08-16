package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FolkLoreee/sapu/internal/app"
)

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

	if err := app.Run(*jpgDir, *rawDir, *hard, *dryRun, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "sapu: %v\n", err)
		os.Exit(1)
	}
}
