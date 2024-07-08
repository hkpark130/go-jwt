FROM golang:1.19

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        git \
        build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src
COPY ./ /go/src

RUN go build main.go
CMD ["go", "run", "main.go"]
