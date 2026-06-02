# build
FROM golang:1.26 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /out/auth ./cmd/auth \
 && CGO_ENABLED=1 go build -o /out/migrator ./cmd/migrator

# run
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /out/ /app/
COPY config ./config
COPY migrations ./migrations
EXPOSE 40049
CMD ["/bin/sh", "-c", "mkdir -p ./storage && STORAGE_PATH=./storage/auth.db MIGRATIONS_PATH=./migrations /app/migrator && CONFIG_PATH=./config/docker.yaml exec /app/auth"]
