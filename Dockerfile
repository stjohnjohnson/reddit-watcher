# Build Step
FROM golang:1.10 AS build

# Prerequisites
WORKDIR /go/src/github.com/stjohnjohnson/reddit-watch
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/screwdriver-cd/gitversion

# Compilation target
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# Add code and install deps
COPY . /go/src/github.com/stjohnjohnson/reddit-watch
RUN dep ensure -vendor-only

# Build the app
RUN go build \
    -o /reddit-watch \
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
COPY --from=build /reddit-watch /usr/bin/

# Specify our launch point
ENTRYPOINT ["/usr/bin/reddit-watch"]