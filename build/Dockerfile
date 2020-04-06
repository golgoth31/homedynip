# Get certificates for https conexions
FROM alpine:latest as certs
RUN apk add -U --no-cache ca-certificates

# build binaries
FROM golang:1.13 as golang-build
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
WORKDIR /go/src/github.com/golgoth31/homedynip
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make artifact_all
