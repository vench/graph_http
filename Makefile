V := @

# build
NAME=graph_http
OUT_DIR=./bin
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`

LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.build=${BUILD}"

.PHONY: build
build:
	$(V)CGO_ENABLED=1 go build ${LDFLAGS} -o ${OUT_DIR}/${NAME} ./

.PHONY: test
test: GO_TEST_FLAGS += -race
test:
	$(V)go test -mod=mod $(GO_TEST_FLAGS) --tags=$(GO_TEST_TAGS) ./...

.PHONY: lint
lint:
	$(V)${OUT_DIR}/golangci-lint run

.PHONY: lint-install
lint-install:
	$(V)curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.50.1
	$(V)${OUT_DIR}/golangci-lint --version