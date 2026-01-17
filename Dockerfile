FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /app/bin

RUN CGO_ENABLED=0 go build -o /app/bin/server ./app-1/cmd/server
RUN CGO_ENABLED=0 go build -o /app/bin/worker ./app-1/cmd/worker
RUN CGO_ENABLED=0 go build -o /app/bin/app2   ./app-2/cmd

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin/server .
COPY --from=builder /app/bin/worker .
COPY --from=builder /app/bin/app2 .
COPY --from=builder /app/app-1/migrations ./migrations
RUN ls -la /app/migrations
CMD ["./server"]
