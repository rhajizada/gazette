FROM golang:1.24-alpine AS builder
WORKDIR /src
ARG VERSION

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
  go build -ldflags="-s -w -X main.Version=${VERSION}" \
  -o /src/bin/scheduler \
  ./cmd/scheduler

FROM alpine:3.18
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /src/bin/scheduler .

ENTRYPOINT ["./scheduler"]
