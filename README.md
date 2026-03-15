# twofold

An opinionated duplicate-finder for local folders.

Twofold scans a directory tree, computes SHA-256 checksums of regular files, groups identical files by hash and size, and can list duplicates or replace duplicates with symlinks or hardlinks (or remove them).

This repository contains a small CLI and tests. Use with care — some operations are destructive.

## Features

- Fast concurrent hashing with a bounded worker pool (default: number of CPU cores).
- Aggregated progress bar when running in a terminal and `-v`/`--verbose` is enabled.
- Safe defaults: dry-run by default — use `--apply` to make changes.
- Hardlink fallback: when hardlinking fails across devices, twofold falls back to copying the file so duplicate content is preserved.
- Unit tests covering synchronous and concurrent scans, and replacement behaviors.

## Install

Build from source:

```sh
go build -o twofold ./...
```

## Usage

Basic listing of duplicates in the current directory:

```sh
twofold -l
```

Specify a path:

```sh
twofold -l ~/Downloads
```

Flags

- `-v`, `--verbose` — show checksumming progress. When concurrency is enabled, progress is shown as a single aggregated bar in a TTY.
- `-l`, `--list` — list duplicates only (dry-run behavior).
- `--symlink` — remove duplicates and create symlinks pointing to the canonical file.
- `--hardlink` — remove duplicates and create hardlinks pointing to the canonical file. If hardlinking fails because files are on different devices, twofold falls back to copying the original file to the duplicate path.
- `--remove` — remove duplicate files without linking them.
- `--apply` — apply destructive changes (without this flag the program runs as a dry-run and will only print what it would do).
- `--workers N` — number of concurrent hashing workers. Default: number of CPU cores. Use `--workers 0` or `--workers 1` to force single-threaded operation.
- `--path PATH` — directory to scan (default: `.`).

Examples

Dry-run listing in Downloads:

```sh
twofold -l --path ~/Downloads
```

Hardlink duplicates (apply changes):

```sh
twofold --hardlink --apply --path /mnt/data/photos
```

Force single-threaded (useful for debugging):

```sh
twofold --workers 1 -l
```

## Safety notes

- The program is potentially destructive when `--apply` is used. Double-check options before running on important data.
- Hardlinks cannot be created across filesystem boundaries. twofold will attempt a fallback copy when that occurs; copied files preserve the original file mode.
- When in doubt, run with `-l` (list) and `--apply` omitted to preview actions.

## Development

- Tests: `go test ./...` (CI runs `go vet` and `go test`).
- Formatting: `gofmt -w .`.

## License

MIT
