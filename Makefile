exec_base = run-go
exec_ext =

# use `go version` as an os-independent way of getting the current os and arch
os_arch = $(word 4, $(shell go version))
os = $(word 1, $(subst /, ,$(os_arch)))
arch = $(word 2, $(subst /, ,$(os_arch)))

ifeq (${os}, windows)
	exec_ext = .exe
endif

build:
	GOOS=${os} GOARCH=${arch} go build -o ${exec_base}-${os}-${arch}${exec_ext} ./
