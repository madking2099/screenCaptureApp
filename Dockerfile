FROM golang:1.21 AS builder
WORKDIR /app
# Copy only go.mod and go.sum first to leverage caching
COPY go.mod go.sum ./
# Download dependencies explicitly for the module
RUN go mod download
# Copy remaining files
COPY main.go .
COPY swagger-ui ./swagger-ui
# Install latest swag (check for latest version on GitHub)
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.2
# Generate OpenAPI 3.0 swagger.json with dynamic host
RUN SERVER_HOST=http://192.168.1.15:1388 swag init -g main.go --output docs --parseDependency --parseInternal --parseDepth 1
RUN cat /app/docs/swagger.json
# Build with CGO disabled and target Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o screenshot-service .

FROM debian:bullseye-slim
WORKDIR /app
RUN apt-get update && apt-get install -y \
    chromium \
    libnss3 \
    libatk1.0-0 \
    libatk-bridge2.0-0 \
    libcups2 \
    libxkbcommon0 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    libgbm1 \
    libpango-1.0-0 \
    libcairo2 \
    libasound2 \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/screenshot-service .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/swagger-ui ./swagger-ui
RUN mkdir -p static
EXPOSE 8000
# Set SERVER_HOST at runtime for UI (optional, can override build-time)
ENV SERVER_HOST=http://192.168.1.15:1388
CMD ["./screenshot-service"]