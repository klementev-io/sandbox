ARG GOLANG_VERSION=1.24
ARG ALPINE_VERSION=3.22

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/sandbox ./cmd/sandbox

FROM alpine:${ALPINE_VERSION}

COPY --from=builder /app/sandbox /sandbox
COPY --from=builder /app/configs/config.yaml /config.yaml

RUN addgroup -S sandbox && adduser -S sandbox -G sandbox
USER sandbox

ENTRYPOINT ["/sandbox", "-c", "/config.yaml"]