# go-fltk-clipboard

A simple clipboard manager.

Features dark/light mode and portrait/landscape mode.

## Screenshots

Coming soon.

## Requirements

Linux requirements:

- x11 (wayland is untested)
- `xclip` or `xsel`

Other platforms are untested but may work. [See compatibility here](https://github.com/atotto/clipboard).

## Installation

This application can be installed via Flatpak:

```bash
# if you do not have flathub added as a remote, please add it first, so that
# the necessary flatpak runtimes can be acquired:
flatpak --user remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo

flatpak --user remote-add --if-not-exists cmcode https://flatpak.cmcode.dev/cmcode.flatpakrepo

flatpak --user install cmcode dev.cmcode.go-fltk-clipboard
```

## Development setup

This repository makes use of `git lfs` for tracking its word dictionaries. Please ensure you have it working.

To build, run

```bash
make build-prod
```

To install to `~/.local/bin/`, run

```bash
make install
```

## Building the flatpak

To build the flatpak, use:

```bash
make flatpak-build-test
```

This will install a local version of the flatpak, and you can run it via

```bash
flatpak --user run dev.cmcode.go-fltk-clipboard
```

Once you're satisfied with it, you can then proceed to release it, assuming the remote repository's mount point is set up correctly:

```bash
# WARNING: This will update the globally available repository!
make flatpak-release
```
