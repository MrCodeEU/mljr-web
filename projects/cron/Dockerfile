# syntax=docker/dockerfile:1

# Build context must be the repo root (shared ui/ packages are needed):
# docker build -f projects/cron/Dockerfile -t mljr/cron .

FROM golang:1.26-bookworm AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make tailwind vendor-js

RUN make build PROJECT=cron

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/bin/cron /app/cron

ENV MLJR_ENV=prod
EXPOSE 8093
ENTRYPOINT ["/app/cron"]
