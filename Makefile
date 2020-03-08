# Generic Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=craig
BINARY_UNIX=$(BINARY_NAME)_unix

# Craig parameters
ENV_FILE='./.env'
CONFIG_FILE='./dev.config.json'
include .env
export

all: test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v main/main.go
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN) ./main
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
#		if [ -f $(ENV_FILE) ]; then
#		fi
		./$(BINARY_NAME) --config-file=$(CONFIG_FILE)
deps:
		go mod download
		go mod verify


# Cross compilation
build-linux:
		# depends on `brew install FiloSottile/musl-cross/musl-cross`
 		# this is currently broken (doesn't compile on macos)
		CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc $(GOBUILD) -o $(BINARY_UNIX) -a -v main/main.go
docker-build:
        #TODO
		#docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
