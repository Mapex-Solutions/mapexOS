# Stage 1: Builder
FROM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

RUN apk add --no-cache git ca-certificates

COPY mapexGoKit/ mapexGoKit/
COPY mapexOS/workspace_go/packages/ mapexOS/workspace_go/packages/
COPY mapexOS/workspace_go/services/mapexIam/ mapexOS/workspace_go/services/mapexIam/

WORKDIR /build/mapexOS/workspace_go/services/mapexIam

ENV GOWORK=off

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -o /app/service \
    ./src/main.go

# Stage 2: Production
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -g 1000 mapexIam && adduser -u 1000 -G mapexIam -s /bin/sh -D mapexIam

RUN mkdir -p /app /var/cache/mapexIam /var/log/mapexIam

COPY --from=builder /app/service /app/service

RUN chown -R mapexIam:mapexIam /app /var/cache/mapexIam /var/log/mapexIam

WORKDIR /app

USER mapexIam

CMD ["/app/service"]
