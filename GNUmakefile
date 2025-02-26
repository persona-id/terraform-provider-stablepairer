default: fmt lint install generate

build:
	go build -v ./...

fmt:
	gofmt -s -w -e .

generate:
	cd tools; go generate ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate
