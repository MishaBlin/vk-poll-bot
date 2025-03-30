FROM golang:1.24 AS builder

RUN apt-get update && apt-get install -y pkg-config libssl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o poll-app ./cmd/mm-polls

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y libssl3 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/poll-app .

EXPOSE 8080

CMD ["./poll-app"]
