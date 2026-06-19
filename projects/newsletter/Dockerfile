# syntax=docker/dockerfile:1

# Build context must be the repo root (shared ui/ packages are needed):
# docker build -f projects/newsletter/Dockerfile -t mljr/newsletter .

FROM golang:1.26-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make tailwind vendor-js

RUN make build PROJECT=newsletter

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/bin/newsletter /app/newsletter

ENV MLJR_ENV=prod
ENV NEWSLETTER_DATA_DIR=/data
VOLUME ["/data"]
EXPOSE 8096
ENTRYPOINT ["/app/newsletter"]
