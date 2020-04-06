# Get certificates for https conexions
FROM alpine:latest as certs
RUN apk add -U --no-cache ca-certificates

# build proto files
# FROM thethingsindustries/protoc:3.1.21 as buildproto
# WORKDIR /src/github.com/golgoth31/homedynip
# ENV GOPATH=/
# RUN mkdir -p proto pkg
# COPY proto proto
# RUN protoc \
#     -Iproto \
#     -I${GOPATH}/src \
#     --go_out=${GOPATH}/src \
#     ./proto/*.proto

# build binary
FROM golang:1.13 as golang-build
ARG GOARCH
ARG GOARM
ENV GOARCH=$GOARCH
ENV GOARM=$GOARM
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org
WORKDIR /go/src/github.com/golgoth31/homedynip
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# COPY --from=buildproto /src/github.com/golgoth31/homedynip/pkg ./pkg
RUN make build_local

FROM scratch
COPY --from=golang-build /go/src/github.com/golgoth31/homedynip/homedynip /homedynip
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8080
ENTRYPOINT ["/homedynip"]
