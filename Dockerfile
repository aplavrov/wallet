FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o wallet ./cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/wallet .

COPY config.env .
COPY internal/storage/db/migrations ./internal/storage/db/migrations

EXPOSE 9000

CMD ["./wallet"]
