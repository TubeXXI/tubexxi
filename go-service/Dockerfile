FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/api ./cmd/api && \
    go build -o bin/worker ./cmd/worker && \
    go build -o bin/seed ./cmd/seed

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    wget \
    tar \
    gzip \
    unzip \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/* && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

COPY --from=builder /app/bin/api ./api
COPY --from=builder /app/bin/worker ./worker
COPY --from=builder /app/bin/seed ./seed

COPY --from=builder /app/database/migrations ./database/migrations

EXPOSE 5002

CMD ["./api"]
