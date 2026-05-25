# Stage 1: Builder
FROM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

RUN apk add --no-cache git ca-certificates

COPY mapexGoKit/ mapexGoKit/
COPY mapexOS/workspace_go/packages/ mapexOS/workspace_go/packages/
COPY mapexOS/workspace_go/services/workflow/ mapexOS/workspace_go/services/workflow/

WORKDIR /build/mapexOS/workspace_go/services/workflow

ENV GOWORK=off

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -o /app/service \
    ./src/main.go

# Stage 2: Production
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -g 1000 workflow && adduser -u 1000 -G workflow -s /bin/sh -D workflow

RUN mkdir -p /app /var/log/workflow /var/cache/workflow

COPY --from=builder /app/service /app/service

RUN chown -R workflow:workflow /app /var/log/workflow /var/cache/workflow

WORKDIR /app

USER workflow

CMD ["/app/service"]
