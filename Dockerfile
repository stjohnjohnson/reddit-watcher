# Build Step
FROM golang:1.10 AS build

# Prerequisites
WORKDIR /go/src/github.com/stjohnjohnson/reddit-watcher
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/screwdriver-cd/gitversion

# Compilation target
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# Add code and install deps
COPY . /go/src/github.com/stjohnjohnson/reddit-watcher
RUN dep ensure -vendor-only

# Test the code
RUN go test ./... -coverprofile coverage.out

# Build the app
RUN go build \
    -o /reddit-watcher \
    -a -installsuffix cgo \
    -ldflags "\
        -X main.version=`gitversion show` \
        -X main.commit=`git rev-parse HEAD` \
        -X main.date=`date -u +"%Y-%m-%dT%H:%M:%SZ"`\
    "

# Executable stage
FROM alpine:3.4

# Ensure we can call HTTPS endpoints
RUN apk add --update ca-certificates && rm -rf /var/cache/apk/*

# Copy binary from build step
COPY --from=build /reddit-watcher /usr/bin/

# Persist data in this directory
VOLUME /config

# Specify our launch point
ENTRYPOINT ["/usr/bin/reddit-watcher"]
