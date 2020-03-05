GOCMD=go
GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
BUILD_LINUX="build/ecalc"
BUILD_WINDOWS="build/ecalc.exe"
BUILD_WINDOWS32="build/ecalc32.exe"

all: test build

test: 
	$(GOTEST) -v ./...

build: build-linux build-windows build-windows32

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -a -o $(BUILD_LINUX)

build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -a -o $(BUILD_WINDOWS)

build-windows32:
	GOOS=windows GOARCH=386 $(GOBUILD) -a -o $(BUILD_WINDOWS32)