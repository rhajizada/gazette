FROM golang:1.24-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
  go build -ldflags="-s -w" \
  -o /src/bin/server \
  ./cmd/server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /src/bin/server .

COPY --from=builder /src/data/sql/migrations data/sql/migrations

ENTRYPOINT ["./server"]
