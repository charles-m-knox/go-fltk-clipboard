.PHONY=build

BUILDDIR=build
VER=0.0.1
BIN=$(BUILDDIR)/go-fltk-clipboard-v$(VER)

build-dev:
	CGO_ENABLED=0 go build -v

mkbuilddir:
	mkdir -p $(BUILDDIR)

build-prod: mkbuilddir
	CGO_ENABLED=0 go build -v -o $(BIN) -ldflags="-w -s -buildid=" -trimpath

test:
	go test -test.v -coverprofile=testcov.out ./... && \
	go tool cover -html=testcov.out

run:
	./$(BIN)

lint:
	golangci-lint run ./...

compress-prod: mkbuilddir
	rm -f $(BIN)-compressed
	upx --best -o ./$(BIN)-compressed $(BIN)

build-mac-arm64: mkbuilddir
	CGO_ENABLED=0 GOARCH=arm64 GOOS=darwin go build -v -o $(BIN)-darwin-arm64 -ldflags="-w -s -buildid=" -trimpath
	rm -f $(BIN)-darwin-arm64.xz
	xz -9 -e -T 12 -vv $(BIN)-darwin-arm64

build-mac-amd64: mkbuilddir
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -v -o $(BIN)-darwin-amd64 -ldflags="-w -s -buildid=" -trimpath
	rm -f $(BIN)-darwin-amd64.xz
	xz -9 -e -T 12 -vv $(BIN)-darwin-amd64

build-win-amd64: mkbuilddir
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -v -o $(BIN)-win-amd64-uncompressed -ldflags="-w -s -buildid=" -trimpath
	rm -f $(BIN)-win-amd64
	upx --best -o ./$(BIN)-win-amd64 $(BIN)-win-amd64-uncompressed

build-linux-arm64: mkbuilddir
	CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -v -o $(BIN)-linux-arm64-uncompressed -ldflags="-w -s -buildid=" -trimpath
	rm -f $(BIN)-linux-arm64
	upx --best -o ./$(BIN)-linux-arm64 $(BIN)-linux-arm64-uncompressed

build-linux-amd64: mkbuilddir
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o $(BIN)-linux-amd64-uncompressed -ldflags="-w -s -buildid=" -trimpath
	rm -f $(BIN)-linux-amd64
	upx --best -o ./$(BIN)-linux-amd64 $(BIN)-linux-amd64-uncompressed

build-all: mkbuilddir build-linux-amd64 build-linux-arm64 build-win-amd64 build-mac-amd64 build-mac-arm64

delete-builds:
	rm $(BUILDDIR)/*
