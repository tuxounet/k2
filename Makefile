APP_NAME := "k2"	
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GO_PATH:=$(shell go env GOPATH)
VERSION := $(if $(CI_COMMIT_TAG),$(CI_COMMIT_TAG),v${GIT_BRANCH})
VERSION_FILE := ./version.yaml


write-version: 
	echo ${VERSION} > ${VERSION_FILE}

run: write-version
	go run ./main.go

build: write-version
	go build  -o ./.out/k2 ./main.go

bump-patch: build 
	echo "Bumping version patch"
	@sh ./scripts/bump-patch.sh
test:
	go test -v ./...

plan:
	go run ./main.go plan  --inventory ./samples/k2.inventory.yaml

apply:
	go run ./main.go apply --inventory ./samples/k2.inventory.yaml

destroy:
	go run ./main.go destroy --inventory ./samples/k2.inventory.yaml