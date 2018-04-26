# Build stage
FROM golang:1.10-alpine3.7
MAINTAINER Hung-Wei Chiu <hwchiu@linkernetworks.com>

WORKDIR /go/src/github.com/linkernetworks/fileserver
COPY src   /go/src/github.com/linkernetworks/fileserver/src
COPY main.go /go/src/github.com/linkernetworks/fileserver
COPY vendor /go/src/github.com/linkernetworks/fileserver/vendor

RUN go install .
ENTRYPOINT /go/bin/fileserver -host localhost -port 33333
