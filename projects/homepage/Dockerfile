# syntax=docker/dockerfile:1

# Build context must be the repo root (mljr-data submodule + shared ui/ packages
# are needed): docker build -f projects/homepage/Dockerfile -t mljr/homepage .

FROM golang:1.26-bookworm AS builder
# debian/glibc base: the prebuilt tailwindcss binary doesn't run on musl/alpine.

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Fetch the tailwind binary + vendored JS/CSS (datastar, altcha, leaflet).
RUN make tailwind vendor-js

# Bake the current mljr-data site data in as the embedded fallback so the
# binary works standalone; HOMEPAGE_DATA_FILE can still point at a mounted
# mljr-data checkout at runtime to override it without a rebuild.
RUN cp mljr-data/generated/site-data.json projects/homepage/data/seed-cache.json

RUN make build PROJECT=homepage

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /src/bin/homepage /app/homepage

ENV MLJR_ENV=prod
EXPOSE 8090
ENTRYPOINT ["/app/homepage"]
