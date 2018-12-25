BIN = maze
DIR = ./cmd/maze

.PHONY: all
all: clean test build

.PHONY: build
build: deps
	go build -o build/$(BIN) $(DIR)

.PHONY: install
install: deps
	go install ./...

.PHONY: deps
deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: cross
cross: crossdeps
	goxz -os=linux,darwin,freebsd,netbsd,windows -arch=386,amd64 -n $(BIN) $(DIR)

.PHONY: crossdeps
crossdeps: deps
	go get github.com/Songmu/goxz/cmd/goxz

.PHONY: test
test: testdeps build
	go test -v $(DIR)...

.PHONY: testdeps
testdeps:
	go get -d -v -t ./...

.PHONY: lint
lint: lintdeps build
	go vet
	golint -set_exit_status $(go list ./... | grep -v /vendor/)

.PHONY: lintdeps
lintdeps:
	go get -d -v -t ./...
	command -v golint >/dev/null || go get -u golang.org/x/lint/golint

.PHONY: clean
clean:
	rm -rf build goxz
	go clean
