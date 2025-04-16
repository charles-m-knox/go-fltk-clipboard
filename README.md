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

First, install Go - this varies based on your OS. Make sure your `$PATH`
includes `go env | grep GOBIN`.

Then, use `go install`:

```bash
CGO_ENABLED=1 go install -ldflags="-w -s -buildid= -X main.version=0.0.5" -trimpath github.com/charles-m-knox/go-fltk-clipboard@v0.0.5
```

## Development setup

To do a proper build, run the following - make sure you have `upx` and `go` installed as well as the listed requirements above:

```bash
make build-prod
```
