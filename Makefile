NAME=krok

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := bin

# VERSION defines the project version for the bundle. 
VERSION ?= 0.0.1

# List the GOOS and GOARCH to build
GO_LDFLAGS_STATIC="-s -w $(CTIMEVAR) -extldflags -static"

.DEFAULT_GOAL := help

##@ Build

binaries: ## Builds binaries for all supported platforms, linux, darwin
	CGO_ENABLED=0 gox \
		-osarch="linux/amd64 linux/arm darwin/amd64" \
		-ldflags=${GO_LDFLAGS_STATIC} \
		-output="$(BUILDDIR)/{{.OS}}/{{.Arch}}/$(NAME)" \
		-tags="netgo" \
		./

bootstrap: ## Installs necessary third party components
	go get github.com/mitchellh/gox

##@ Testing

test-db: ## Generates a Test database to use for local development of Krok
	docker run -d \
		--rm \
		-e POSTGRES_USER=krok \
		-e POSTGRES_PASSWORD=password123 \
		-v `pwd`/dbinit:/docker-entrypoint-initdb.d \
		-p 5432:5432 \
		--name krok-test-db \
		postgres:13.1-alpine

rm-test-db: ## Removes the created test database.
	docker rm -f krok-test-db

test: lint ## Lints Krok then runs all tests
	go test -count=1 ./...

clean: ## Runs go clean
	go clean -i

lint: ## Runs golangci-lint on Krok
	golangci-lint run ./...

##@ Docker

docker_image: ## Creates a docker image for Krok. Requires `image` and `version` variables on command line
	docker build -t $(image):$(version) .

##@ Utilities

generate_mocks: ## Generates all mocks for all defined interfaces for testing purposes.
	go build -o pkg/krok/providers/interfaces generate_interface_mocks/main.go && cd pkg/krok/providers && ./interfaces && rm ./interfaces

generate-protoc: ## Generate the protobuf files
	buf generate

protoc-update: ## Update the protobuf modules
	buf mod update

lint-protoc: ## lint the protocol files
	buf lint
	
help:  ## Display this help. Thanks to https://www.thapaliya.com/en/writings/well-documented-makefiles/
ifeq ($(OS),Windows_NT)
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-40s %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
else
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-40s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
endif