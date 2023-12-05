# use `go version` as an os-independent way of getting the current os and arch
OS_ARCH = $(word 4, $(shell go version))
OS = $(word 1, $(subst /, ,$(OS_ARCH)))
ARCH = $(word 2, $(subst /, ,$(OS_ARCH)))

BINARY_NAME = run-go

build-darwin:
ifeq (${ARCH}, amd64)
	GOARCH=${ARCH} GOOS=darwin go build -o ${BINARY_NAME}-darwin-${ARCH} ./
else ifeq (${ARCH}, arm64)
	GOARCH=${ARCH} GOOS=darwin go build -o ${BINARY_NAME}-darwin-${ARCH} ./
endif

build-linux:
ifeq (${ARCH}, amd64)
	GOARCH=${ARCH} GOOS=linux go build -o ${BINARY_NAME}-linux-${ARCH} ./
else ifeq (${ARCH}, arm64)
	GOARCH=${ARCH} GOOS=linux go build -o ${BINARY_NAME}-linux-${ARCH} ./
endif

build-windows:
ifeq (${ARCH}, amd64)
	GOARCH=${ARCH} GOOS=windows go build -o ${BINARY_NAME}-windows-${ARCH}.exe ./
endif
