NAME=krok

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := bin

# List the GOOS and GOARCH to build
GO_LDFLAGS_STATIC="-s -w $(CTIMEVAR) -extldflags -static"

.DEFAULT_GOAL := binaries

.PHONY: binaries
binaries:
	CGO_ENABLED=0 gox \
		-osarch="linux/amd64 linux/arm darwin/amd64" \
		-ldflags=${GO_LDFLAGS_STATIC} \
		-output="$(BUILDDIR)/{{.OS}}/{{.Arch}}/$(NAME)" \
		-tags="netgo" \
		./...

.PHONY: bootstrap
bootstrap:
	go get github.com/mitchellh/gox

.PHONY: test-db
test-db:
	docker run -d \
		--rm \
		-e POSTGRES_USER=krok \
		-e POSTGRES_PASSWORD=password123 \
		-v `pwd`/dbinit:/docker-entrypoint-initdb.d \
		-p 5432:5432 \
		--name krok-test-db \
		postgres:13.1-alpine

.PHONY: rm-test-db
rm-test-db:
	docker rm -f krok-test-db

# Check if we are in circleci. If yes, start a postgres docker instance.
.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean -i

lint:
	golint ./...

.PHONY: run
run:
	go run cmd/root.go

.PHONY: start-https
start-https:
	go run cmd/root.go --server-key-path ./certs/key.pem --server-crt-path ./certs/cert.pem

docker_image:
	docker build -t $(image):$(version) .