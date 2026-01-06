FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . .

RUN go build -o geo ./cmd/server/main.go
RUN go build -o geo ./cmd/worker/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/worker .