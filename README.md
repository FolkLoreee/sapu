# sapu

A small CLI for cleaning up leftover RAW files after a photo cull.

When you shoot RAW+JPEG and delete the rejects as JPEGs, the orphaned RAW files
pile up. `sapu` removes every RAW file whose stem has no matching JPEG, keeping
the ones you actually kept.

## Install

```sh
go install github.com/FolkLoreee/sapu@latest
```

## Usage

```sh
sapu --jpg /path/to/jpegs --raw /path/to/raws
```

| Flag | Default | Description |
|------|---------|-------------|
| `--jpg` | — | (required) Directory containing JPEG files. |
| `--raw` | — | (required) Directory containing RAW files. |
| `--hard` | `false` | Permanently delete instead of moving to the macOS Trash. |
| `--dry-run` | `false` | Print what would be removed without touching anything. |

By default files are moved to the macOS Trash via Finder, so a mistake is
recoverable. Use `--dry-run` to preview first:

```sh
sapu --jpg /path/to/jpegs --raw /path/to/raws --dry-run
sapu --jpg /path/to/jpegs --raw /path/to/raws --hard
```

## Matching rules

- A RAW file is kept when its **stem** matches a JPEG's stem. Stems are the
  filename with the final extension removed, compared case-sensitively.
- Extensions are matched case-insensitively.
- JPEG: `.jpg`, `.jpeg`
- RAW: `.cr2`, `.cr3`, `.nef`, `.nrw`, `.arw`, `.dng`, `.raf`, `.orf`, `.rw2`,
  `.pef`, `.srw`
- If one stem matches, **all** RAW files with that stem are kept.
- Only the top level of each directory is scanned; subdirectories, dotfiles, and
  symlinks are ignored.

## Example

```
jpegs/          raws/
  img1.jpg        img1.cr2
  img5.jpg        img2.cr2
                  img3.cr2
                  img4.cr2
                  img5.cr2
```

```
$ sapu --jpg jpegs --raw raws --dry-run
[DRY-RUN] would remove: raws/img2.cr2
[DRY-RUN] would remove: raws/img3.cr2
[DRY-RUN] would remove: raws/img4.cr2
scanned 5 raw files: 2 kept, 3 removed
```

## Build

```sh
go build
```

## Layout

```
main.go                 flag parsing and command routing
internal/app            orchestration and output
internal/filetype       extension classification and stem derivation
internal/match          keep/delete classification
internal/scan           directory listing
internal/remove         permanent and Trash deletion
internal/testutil       shared test helpers
```
