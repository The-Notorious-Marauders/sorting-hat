FROM golang:1.21-alpine AS builder
ARG MODULE=rest
ARG BUILD_VERSION
ARG BUILD_COMMIT_HASH
ARG BUILD_TIME
ARG BS_PKG=github.com/The-Notorious-Marauders/sorting-hat/${MODULE}/bootstrap
RUN apk add --no-cache git
COPY ./src /go/src
WORKDIR /go/src/${MODULE}
RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./application \
    -ldflags="-X '${BS_PKG}.Version=${BUILD_VERSION}' -X '${BS_PKG}.CommitHash=${BUILD_COMMIT_HASH}' -X '${BS_PKG}.BuildTime=${BUILD_TIME}'"

FROM alpine:3.15
ARG MODULE=rest
COPY --from=builder /go/src/${MODULE}/config /app/config
COPY --from=builder /go/src/${MODULE}/application /app
EXPOSE 3700
WORKDIR /app
ENTRYPOINT ["./application"]
