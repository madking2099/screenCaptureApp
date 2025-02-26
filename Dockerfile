FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
RUN go mod tidy  # Populate go.mod with indirect dependencies
RUN go get -d -v ./...  # Download all dependencies explicitly
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