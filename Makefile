GO111MODULE := on
export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)

.PHONY: rcid
rcid:
	GOOS=darwin \
	CGO_ENABLED=0 \
	go build \
	-v \
	-o build/rcid \
	${PWD}/main.go

.PHONY: rcidcli
rcidcli:
	GOOS=darwin \
	CGO_ENABLED=0 \
	go build \
	-v \
	-o build/rcidcli \
	${PWD}/client/rci.go

.PHONY: protoc
protoc:
	docker run --rm \
	-v "${PWD}/proto:/usr/local/src" \
	-v "${PWD}/pb:/usr/local/gen" \
	shoma/protoc

.PHONY: bootstrap-lint-tools
install-lint-tools:
	go get github.com/alecthomas/gometalinter && \
	gometalinter --install --update

.PHONY: fmt
fmt:
	@find . -iname "*.go" -not -path "./vendor/**" | xargs gofmt -s -w

.PHONY: clean
clean:
	rm -fr build/*
