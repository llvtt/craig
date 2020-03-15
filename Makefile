# Generic Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=craig
BINARY_NAME_LAMBDA=main
LAMBDA_DIR=./lambda/main
TERRAFORM_ENV=prod
TERRAFORM_DIR=./terraform/environments/${TERRAFORM_ENV}
LAMBDA_BUILD_FLAGS=-ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o

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
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc $(GOBUILD) -o $(BINARY_NAME) -a -v main/main.go

build-lambda: clean-lambda
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc $(GOBUILD) -o $(LAMBDA_DIR)/$(BINARY_NAME_LAMBDA) $(LAMBDA_BUILD_FLAGS) -a -v $(LAMBDA_DIR)
#build-lambda: clean-lambda
#	cd $(LAMBDA_DIR) && $(GOBUILD) -o $(BINARY_NAME_LAMBDA) -a -v main.go

deploy-plan:  build-lambda
	cd $(TERRAFORM_DIR) && terraform init && terraform plan

deploy: build-lambda
	cd $(TERRAFORM_DIR) && terraform init && terraform apply -auto-approve