# twofold [![Build Status](https://travis-ci.org/Urethramancer/twofold.svg?branch=master)](https://travis-ci.org/Urethramancer/ftwofold)
An exercise in duplicate-finding.

## Purpose
Find duplicates within a folder and list, link or remove them. Slightly dangerous, no backsies if you screw up.

## Usage
The simplest invocation is to list duplicates in the current directory. Example:

```sh
$ twofold -l
Deep-scanning /Users/orb/go/twofold
```

You can specify a path as its sole standalone argument:

```sh
$ twofold -l ~/Downloads
Deep-scanning /Users/orb/Downloads
```

All duplicates found will be listed, grouped by hash, then the program exits.

Use the `-v` flag to also display checksumming progress while it traverses the path.

The slightly more destructive flags are:

- `--symlink` to remove and symlink duplicates from the first file in the set
- `--hardlink` to remove and link (a.k.a. hardlink) duplicates to the inode of the first file in the set (the most convenient, should the original move or be deleted)
- `--remove` to simply remove the duplicates (the most destructive option)

## LICENCE
MIT.

### TODO
- Maybe add date-sorting to make it possible to keep the oldest or newest of each set of duplicates.
- Search multiple supplied directories, finding duplicates between all of them.
