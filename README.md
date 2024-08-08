# go-fltk-clipboard

A simple clipboard manager.

## Requirements

Linux requirements:

- x11 (wayland is untested)
- `xclip` or `xsel`

Other platforms are untested but may work. [See compatibility here](https://github.com/atotto/clipboard).

## Compiling for other targets

I never got this fully working, but here are some notes.

To compile for other platforms, install the following arch linux packages:

- gcc-go
- clang
- aarch64-linux-gnu-gcc

Then, run something like:

```bash
GOARCH=arm64 GOOS=linux CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ go build -v
```
