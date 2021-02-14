#!/bin/bash

TARGETS=(\
         slack-events\
             scraper\
    )
LINUXFLAGS=GOOS=linux GOARCH=amd64
TERRAFORM_DIR=./terraform

build() {
    for target in ${TARGETS[@]}; do
        echo building $target
        go build -o $target cmd/${target}/*.go
    done
}

build_linux() {
    GOOS=linux GOARCH=amd64 build
}

deps() {
    go mod download
    go mod verify
}

clean() {
    go clean
    rm -f ${TARGETS[@]}
}

deploy_plan() {
    cd $TERRAFORM_DIR
    terraform init
    terraform plan
}

deploy() {
    cd $TERRAFORM_DIR
    terraform init
    terraform apply -auto-approve
}

help() {
    cat <<EOF
USAGE:

        $0 build        - Build for linux
        $0 build-dev    - Build locally
        $0 clean        - Clean
        $0 deploy-plan  - Output plan for deploy
        $0 deploy       - Deploy
        $0 deps         - Download and verify dependencies
EOF
}

case "$1" in
    build-dev)
        build $*
        ;;
    build)
        build_linux $*
        ;;
    clean)
        clean $*
        ;;
    deploy-plan)
        deploy_plan $*
        ;;
    deploy)
        deploy $*
        ;;
    help)
        help
        ;;
    *)
        help
        exit 1
esac
