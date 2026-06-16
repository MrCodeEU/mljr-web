# syntax=docker/dockerfile:1

# Build context must be the repo root (shared ui/ packages are needed):
# docker build -f projects/regex/Dockerfile -t mljr/regex .

FROM golang:1.26-bookworm AS builder
# debian/glibc base: the prebuilt tailwindcss binary doesn't run on musl/alpine.

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Fetch the tailwind binary + vendored JS (datastar).
RUN make tailwind vendor-js

RUN make build PROJECT=regex

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/bin/regex /app/regex

ENV MLJR_ENV=prod
EXPOSE 8092
ENTRYPOINT ["/app/regex"]
