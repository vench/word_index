V := @

# Build
OUT_DIR = ./bin
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`

LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.build=${BUILD}"

.PHONY: vendor
vendor:
	$(V)go mod tidy -compat=1.17
	$(V)go mod vendor

.PHONY: build
build:
	$(V)CGO_ENABLED=1 go build ${LDFLAGS} -o ${OUT_DIR}/build ./...

.PHONY: test
test: GO_TEST_FLAGS += -race
test:
	$(V)go test -mod=vendor $(GO_TEST_FLAGS) --tags=$(GO_TEST_TAGS) ./...

.PHONY: lint
lint:
	$(V)./bin/golangci-lint run

.PHONY: lint-install
lint-install:
	$(V)wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.44.0
	$(v)./bin/golangci-lint --version


.PHONY: generate
generate:
	$(V)go generate  ./...

