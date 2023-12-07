# Compile RunGo for Linux/Windows AMD64
FROM ubuntu:latest AS linux-windows-amd64

WORKDIR /app

# Update and install packages
RUN apt-get update
RUN apt-get install -y -q curl libgl1-mesa-dev xorg-dev

# Install Go
RUN curl -OL https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
RUN tar -C /usr/local -xvf go1.21.0.linux-amd64.tar.gz
ENV PATH "$PATH:/usr/local/go/bin"
RUN echo $(go version)

# Install Zig
RUN curl -OL https://ziglang.org/download/0.11.0/zig-linux-x86_64-0.11.0.tar.xz
RUN mkdir /usr/local/zig && tar -C /usr/local/zig -xvf zig-linux-x86_64-0.11.0.tar.xz --strip-components 1
ENV PATH "$PATH:/usr/local/zig"
RUN echo $(zig version)

# Switch to local volume
WORKDIR /run-go

# Build for linux-amd64
RUN mkdir -p dist/linux-amd64
RUN CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64 \
	CC="zig cc -target x86_64-linux-gnu -isystem /usr/include -L/usr/lib/x86_64-linux-gnu" \
	CXX="zig c++ -target x86_64-linux-gnu -isystem /usr/include -L/usr/lib/x86_64-linux-gnu" \
	go build -trimpath -o dist/linux-amd64 .

# Build for windows-amd64
RUN mkdir -p dist/windows-amd64
RUN CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	CC="zig cc -target x86_64-windows-gnu" \
	CXX="zig c++ -target x86_64-windows-gnu" \
	go build -trimpath -ldflags="-H=windowsgui" -o dist/windows-amd64 .

