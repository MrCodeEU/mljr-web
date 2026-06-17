# syntax=docker/dockerfile:1

# Build context must be the repo root (shared ui/ packages are needed):
# docker build -f projects/codec/Dockerfile -t mljr/codec .

FROM golang:1.26-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make tailwind vendor-js

RUN make build PROJECT=codec

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/bin/codec /app/codec

ENV MLJR_ENV=prod
EXPOSE 8094
ENTRYPOINT ["/app/codec"]
