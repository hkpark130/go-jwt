FROM golang:1.17.0-alpine3.14

RUN apk update && apk add git
RUN apk add --update alpine-sdk

WORKDIR /go/src
COPY ./ /go/src

RUN go build main.go
CMD ["go", "run", "main.go"]
