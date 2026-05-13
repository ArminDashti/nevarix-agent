# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build
RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETARCH=amd64
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -trimpath \
    -ldflags="-w -s -X main.buildVersion=${VERSION}" \
    -o /out/nevarix-agent \
    ./cmd/agent-server

FROM alpine:3.19

RUN apk add --no-cache ca-certificates \
    && adduser -D -H -u 65532 appuser \
    && mkdir -p /home/.nevarix-server \
    && chown appuser:appuser /home/.nevarix-server

WORKDIR /app

COPY --from=build --chown=appuser:appuser /out/nevarix-agent /app/nevarix-agent

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/nevarix-agent"]
CMD ["agent"]
