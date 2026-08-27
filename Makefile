.PHONY: all fmt fmt-check lint lintmax gosec govulncheck test test-nocgo race vet build ci-check goreleaser tag-major tag-minor tag-patch release install demo

BINARY := optkit
GORELEASER_ARGS ?= --skip=sign --snapshot --clean
GORELEASER_TARGET ?= --single-target

all: build

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './tmp/*' -not -path './ttmp/*')

fmt-check:
	@files="$$(gofmt -l $$(find . -name '*.go' -not -path './tmp/*' -not -path './ttmp/*'))"; \
	if [ -n "$$files" ]; then echo "$$files"; exit 1; fi

lint:
	GOWORK=off golangci-lint run -v

lintmax:
	GOWORK=off golangci-lint run -v --max-same-issues=100

gosec:
	GOWORK=off go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -exclude-generated -exclude=G101,G204,G301,G304,G306 -exclude-dir=.history ./...

govulncheck:
	GOWORK=off go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

test:
	GOWORK=off CGO_ENABLED=1 go test ./... -count=1

test-nocgo:
	GOWORK=off CGO_ENABLED=0 go test ./... -count=1

race:
	GOWORK=off CGO_ENABLED=1 go test -race ./... -count=1

vet:
	GOWORK=off CGO_ENABLED=1 go vet ./...

build:
	GOWORK=off CGO_ENABLED=1 go build ./...

ci-check: fmt-check vet test test-nocgo build

goreleaser:
	GOWORK=off goreleaser release $(GORELEASER_ARGS) $(GORELEASER_TARGET)

tag-major:
	git tag $$(svu major)

tag-minor:
	git tag $$(svu minor)

tag-patch:
	git tag $$(svu patch)

release:
	git push origin --tags
	GOWORK=off GOPROXY=proxy.golang.org go list -m github.com/go-go-golems/optkit@$$(svu current)

install:
	GOWORK=off CGO_ENABLED=1 go build -o "$$(go env GOPATH)/bin/$(BINARY)" ./cmd/optkit

demo:
	GOWORK=off CGO_ENABLED=1 go run ./cmd/optkit demo --store ./tmp/demo --reset
