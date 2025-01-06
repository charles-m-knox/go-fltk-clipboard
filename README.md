# go-fltk-clipboard

A simple clipboard manager.

Features dark/light mode and portrait/landscape mode and filtering results.

## Screenshots

![Dark mode, portrait orientation](./docs/screenshot-dark-portrait.png)

![Light mode, landscape orientation](./docs/screenshot-light-landscape.png)

## Requirements

Linux requirements:

- x11 (wayland is untested)
- `xclip` or `xsel`

Other platforms are untested but may work. [See compatibility here](https://github.com/atotto/clipboard).

## Installation

Go to the [releases page](https://github.com/charles-m-knox/go-fltk-clipboard/releases) and download the latest version there. Place it anywhere in your `$PATH` and you're good to go.

## Development setup

To build, run the following - make sure you have `upx` and `go` installed as well as the listed requirements above:

```bash
make build-prod
```

To install to `~/.local/bin/`, run

```bash
make install
```

## Flatpak note

It is possible to build a flatpak distribution of this application, but I don't currently have time to deal with the extra overhead, so the only recommended installation method is by building it yourself or downloading a precompiled binary from the releases page (if available).
