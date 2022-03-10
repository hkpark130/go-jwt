FROM golang:1.17.0-alpine3.14

RUN apk update && apk add git

WORKDIR /go/src

RUN go get golang.org/x/tools/cmd/godoc
