PKGS?=$(shell go list ./... | grep -v /vendor/ | tr '\n' ' ')
PREFIX?=$(shell pwd)

binary:
	@echo "+ $@"
	@go build -o ./bin/authz-plugin

clean:
	@echo "+ $@"
	@rm -rf "${PREFIX}/bin/"

# run all lint functionality - excludes Godep directory, vendoring, binaries, python tests, and git files
# Shamelessly inspired from notary
lint:
	@echo "+ $@: go vet, go fmt, golint, misspell, ineffassign"
	# gofmt
	@test -z "$$(gofmt -s -l .| grep -v vendor/ | tee /dev/stderr)"
	# govet
ifeq ($(shell uname -s), Darwin)
	@test -z "$(shell find . -iname *test*.go | grep -v _test.go | grep -v vendor | xargs echo "This file should end with '_test':"  | tee /dev/stderr)"
else
	@test -z "$(shell find . -iname *test*.go | grep -v _test.go | grep -v vendor | xargs -r echo "This file should end with '_test':"  | tee /dev/stderr)"
endif
	@test -z "$$(go tool vet -printf=false . 2>&1 | grep -v vendor/ | tee /dev/stderr)"
	# golint - requires that the following be run first:
	#	go get -u github.com/golang/lint/golint
	# golint
	@test -z "$(shell find . -type f -name "*.go" -not -path "./vendor/*" -exec golint {} \; | tee /dev/stderr)"
	# misspell - requires that the following be run first:
	#    go get -u github.com/client9/misspell/cmd/misspell
	@test -z "$$(find . -type f | grep -v vendor/ | grep -v bin/ | grep -v .git/ | grep -v \.pdf | xargs misspell | tee /dev/stderr)"
	# ineffassign - requires that the following be run first:
	#    go get -u github.com/gordonklaus/ineffassign
	@test -z "$(shell find . -type f -name "*.go" -not -path "./vendor/*" -exec ineffassign {} \; | tee /dev/stderr)"

test: lint
	@echo
	go test $(PKGS)