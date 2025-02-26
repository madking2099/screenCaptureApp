FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod .
RUN go mod tidy
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o screenshot-service main.go

FROM debian:bullseye-slim
WORKDIR /app
RUN apt-get update && apt-get install -y \
    chromium \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/screenshot-service .
RUN mkdir -p static
EXPOSE 8000
CMD ["./screenshot-service"]