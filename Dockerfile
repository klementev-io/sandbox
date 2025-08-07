ARG GOLANG_VERSION="1.24"
ARG ALPINE_VERSION="3.22"

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /tmp/sandbox ./cmd/sandbox

FROM alpine:${ALPINE_VERSION}

RUN addgroup -S sandbox && adduser -S sandbox -G sandbox && \
    mkdir -p /etc/sandbox

COPY --from=builder /tmp/sandbox /usr/local/bin/sandbox

USER sandbox

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/sandbox"]
CMD ["--config", "/etc/sandbox/config.yaml"]