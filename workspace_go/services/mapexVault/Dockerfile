# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates

COPY mapexGoKit/ mapexGoKit/
COPY mapexOS/workspace_go/packages/ mapexOS/workspace_go/packages/
COPY mapexOS/workspace_go/services/mapexVault/ mapexOS/workspace_go/services/mapexVault/

WORKDIR /build/mapexOS/workspace_go/services/mapexVault

ENV GOWORK=off

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/service \
    ./src/main.go

# Stage 2: Production
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -g 1000 mapexVault && adduser -u 1000 -G mapexVault -s /bin/sh -D mapexVault

RUN mkdir -p /app /var/log/mapexVault

COPY --from=builder /app/service /app/service

RUN chown -R mapexVault:mapexVault /app /var/log/mapexVault

WORKDIR /app

USER mapexVault

CMD ["/app/service"]
