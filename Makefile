# Generic Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=craig
BINARY_NAME_LAMBDA=main
LAMBDA_DIR=./lambda/main

# Craig parameters
CONFIG_FILE='./dev.config.json'
export

all: test build
build: clean
		$(GOBUILD) -o $(BINARY_NAME) -v main/main.go
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN) ./main
		rm -f $(BINARY_NAME)
clean-lambda:
		$(GOCLEAN) ./main
		rm -f $(LAMBDA_DIR)/$(BINARY_NAME)
run:
		./$(BINARY_NAME) --config-file=$(CONFIG_FILE)
deps:
		go mod download
		go mod verify


# Cross compilation
build-linux: clean
		# depends on `brew install FiloSottile/musl-cross/musl-cross`
 		# this is currently broken (doesn't compile on macos)
		CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc $(GOBUILD) -o $(BINARY_NAME) -a -v main/main.go

build-lambda: clean-lambda
		# depends on `brew install FiloSottile/musl-cross/musl-cross`
 		# this is currently broken (doesn't compile on macos)
		CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc $(GOBUILD) -o $(LAMBDA_DIR)/$(BINARY_NAME_LAMBDA) -a -v $(LAMBDA_DIR)/main.go

docker-build:
        #TODO
		#docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_NAME)" -v
