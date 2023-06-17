# Build Geth in a stock Go builder container
FROM golang:1.20.5-alpine3.18 as builder

ENV HOME_DIR=/usr/app

RUN apk add --no-cache gcc musl-dev linux-headers git

# Get dependencies - will also be cached if we won't change go.mod/go.sum
WORKDIR $HOME_DIR

COPY go.mod .
COPY go.sum .
RUN go mod download

ADD . .
RUN go build -v ./cmd/compressfile

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /usr/app/compressfile /usr/local/bin/compressfile
