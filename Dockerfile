FROM golang:1.19.3-alpine3.16

RUN go install github.com/magefile/mage@v1.14.0

WORKDIR /go/src

COPY go.mod .
COPY go.sum .

RUN go mod download
