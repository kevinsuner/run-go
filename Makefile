.PHONY: help

help:
	# Thanks rici :) -> https://stackoverflow.com/users/1566221/rici
	@awk '/^#/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:]_-]+:/{print substr($$1,1,index($$1,":")),c}1{c=0}' $(MAKEFILE_LIST) | column -s: -t

# Build and run the application
build-run:
	@echo ">>> Starting app"
	@go mod tidy -go=1.21 && go build -ldflags="-s -w" . && ./run-go

# Build the application for Linux, MacOS and Windows
build: build-linux build-darwin build-windows

# Build the application for Linux (amd64, arm64)
build-linux:
	@echo ">>> Building Go binaries for Linux (amd64, arm64)"
	@fyne-cross linux -arch=amd64,arm64

# Build the application for MacOS (amd64, arm64)
build-darwin:
	@echo ">>> Downloading MacOSX SDK 11.3"
	@curl -OL https://github.com/phracker/MacOSX-SDKs/releases/download/11.3/MacOSX11.3.sdk.tar.xz
	tar -xvf MacOSX11.3.sdk.tar.xz
	@echo ">>> Building Go binaries for Darwin (amd64, arm64)"
	@fyne-cross darwin -macosx-version-min="11.3" -macosx-sdk-path="./MacOSX11.3.sdk" -arch=amd64,arm64
	@echo ">>> Cleaning up before finishing"
	rm MacOSX11.3.sdk.tar.xz
	rm -r MacOSX11.3.sdk 

# Build the application for Windows (amd64)
build-windows:
	@echo ">>> Building Go binaries for Windows (amd64)"
	@fyne-cross windows -arch=amd64

# Install development tools
install-tools:
	@echo ">>> Installing development tools"
	@go install github.com/fyne-io/fyne-cross@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run golangci-lint
lint:
	@echo "Running golangci linters"
	@golangci-lint run
