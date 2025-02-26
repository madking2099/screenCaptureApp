FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
RUN go mod tidy
RUN go get -d -v ./...
# Install swag CLI and generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o screenshot-service main.go

FROM debian:bullseye-slim
WORKDIR /app
RUN apt-get update && apt-get install -y \
    chromium \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/screenshot-service .
COPY --from=builder /app/docs ./docs  # Copy generated Swagger files
RUN mkdir -p static
EXPOSE 8000
CMD ["./screenshot-service"]