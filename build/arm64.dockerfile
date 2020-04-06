FROM alpine:latest as certs
RUN apk add -U --no-cache ca-certificates

FROM scratch AS arm64
COPY homedynip-linux-arm64 /homedynip
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8080
ENTRYPOINT ["/homedynip"]
