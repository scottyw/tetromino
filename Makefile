GO_INSTALLED := $(shell which go)
PORTAUDIO_INSTALLED := $(shell pkg-config --cflags  -- portaudio-2.0 2> /dev/null)

all: deps build test

deps:
ifndef GO_INSTALLED
	$(error "go is not installed, please run 'brew install go' for MacOS or equivalent for your operating system")
endif
ifndef PORTAUDIO_INSTALLED
	$(error "tetromino uses portaudio but it is not installed, please run 'brew install pkg-config portaudio' for MacOS or equivalent for your operating system")
endif

build: deps
	@rm -rf bin
	@mkdir -p bin
	@go build -o bin/tetromino cmd/tetromino/*.go
	@echo "💚 Binaries can be found in 'bin' dir"

test: deps
	@go test ./...
	@echo "💚 All tests completed"
