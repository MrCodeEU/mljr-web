# syntax=docker/dockerfile:1

# Build context must be the repo root (shared ui/ packages are needed):
# docker build -f projects/showcase/Dockerfile -t mljr/showcase .

FROM golang:1.26-bookworm AS builder
# debian/glibc base: the prebuilt tailwindcss binary doesn't run on musl/alpine.

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Fetch the tailwind binary + vendored JS/CSS (datastar, altcha, leaflet).
RUN make tailwind vendor-js

RUN make build PROJECT=showcase

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/bin/showcase /app/showcase

EXPOSE 8091
ENTRYPOINT ["/app/showcase"]
