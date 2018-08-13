SRC_FILES := $(shell find . -name '*.go' | grep -v vendor)
GOOS := $(shell go env | grep GOOS | sed 's/GOOS=//')
GOARCH := $(shell go env | grep GOARCH | sed 's/GOARCH=//')
INSTALL_DIR := /usr/local/bin

build: clean cli

clean:
	rm -rf ./dist

lint:
	gometalinter $(SRC_FILES)

test:
	go test -v ./...

dist:
	mkdir -p dist

task: dist/task
cli: dist/task
dist/task: dist
	go build -o ./dist/task ./cmd/task

install: package install-artifacts

install-snapshot: snapshot install-artifacts

install-artifacts:
	# TODO: Add docs and man pages here once written
	cp dist/${GOOS}_${GOARCH}/task /usr/local/bin/

snapshot: clean
	goreleaser release --skip-publish --snapshot

package: clean
	goreleaser release --skip-publish	

tag-%:
	git tag v$*
	git push --tags

publish:
	goreleaser

release-%: clean tag-% publish
