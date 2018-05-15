install:
	# CI dependencies
	go get -u github.com/golang/dep/cmd/dep
	go get -u gopkg.in/alecthomas/gometalinter.v2
	go get -u github.com/mattn/goveralls
	# Update linter
	gometalinter.v2 --install
	# App dependencies
	dep ensure -vendor-only

test:
	# Lint checks
	gometalinter.v2 ./... --vendor --deadline 2m
	# Tests
	go test ./... -v -covermode=count -coverprofile=coverage.out
	# Code coverage
	@goveralls -coverprofile=coverage.out -service=circle-ci -repotoken ${COVERALLS_TOKEN}

build:
	docker build . -t reddit-watcher:latest

run:
	docker run --rm -it -v `pwd`/config:/config reddit-watcher:latest --token ${TOKEN}

bump:
	# Ensure we have gitversion
	go get -u github.com/screwdriver-cd/gitversion
	# Bump version
	gitversion bump auto
	# Push new tags
	git push origin --tags
