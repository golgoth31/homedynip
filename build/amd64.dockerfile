FROM alpine:latest as certs
RUN apk add -U --no-cache ca-certificates

FROM --platform=linux/amd64 scratch AS amd64
COPY homedynip-linux-amd64 /homedynip
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8080
ENTRYPOINT ["/homedynip"]
