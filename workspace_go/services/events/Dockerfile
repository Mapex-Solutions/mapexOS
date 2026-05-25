# Stage 1: Builder
FROM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

RUN apk add --no-cache git ca-certificates

COPY mapexGoKit/ mapexGoKit/
COPY mapexOS/workspace_go/packages/ mapexOS/workspace_go/packages/
COPY mapexOS/workspace_go/services/events/ mapexOS/workspace_go/services/events/

WORKDIR /build/mapexOS/workspace_go/services/events

ENV GOWORK=off

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -o /app/service \
    ./src/main.go

# Stage 2: Production
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -g 1000 mapexos && adduser -u 1000 -G mapexos -s /bin/sh -D mapexos

RUN mkdir -p /app /var/cache/mapexos /var/log/mapexos

COPY --from=builder /app/service /app/service

RUN chown -R mapexos:mapexos /app /var/cache/mapexos /var/log/mapexos

WORKDIR /app

USER mapexos

CMD ["/app/service"]
