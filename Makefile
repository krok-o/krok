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
		./

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
	go test -count=1 ./...

.PHONY: clean
clean:
	go clean -i

lint:
	golint ./...

.PHONY: run
run:
	go run main.go

docker_image:
	docker build -t $(image):$(version) .

generate_mocks:
	go build -o pkg/krok/providers/interfaces generate_interface_mocks/main.go && cd pkg/krok/providers && ./interfaces && rm ./interfaces

.PHONY: swagger
swagger:
	swagger generate spec -o ./swagger/swagger.yaml -c server -c handlers -c main -c docs -c models --scan-models

.PHONY: swagger-server
swagger-serve:
	swagger serve -F=swagger ./swagger/swagger.yaml