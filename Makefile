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
	# @TODO uncomment after PR #1 is merged
	# gometalinter.v2 ./... --vendor
	go test ./... -v -covermode=count -coverprofile=coverage.out
	@goveralls -coverprofile=coverage.out -service=circle-ci -repotoken ${COVERALLS_TOKEN}

build:
	docker build . -t reddit-watch:latest

run:
	docker run --rm -it -v `pwd`/config:/config reddit-watch:latest --token ${TOKEN}
