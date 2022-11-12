FROM golang:1.19-alpine as builder
COPY ./ /build

RUN cd /build && \
    CGO_ENABLED=0 go build -gcflags "all=-N -l" -ldflags '-extldflags "-static"' -o main

FROM alpine as app
COPY --from=builder /build/main /app/

WORKDIR /app

EXPOSE 2345 8080