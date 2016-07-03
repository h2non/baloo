default: all

all: test vet lint

test:
	go test -v -race ./...
	go test -v -race ./_examples/*/

fmt:
	gofmt -s -d ./...

lint:
	golint ./...

vet:
	go vet ./...

sloc:
	wc -l */**.go

update:
	go get -u ./...

.PHONY: lint test
