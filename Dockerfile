FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
COPY swagger-ui ./swagger-ui
RUN go mod tidy
RUN go get -d -v ./...
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go && cat /app/docs/swagger.json
RUN CGO_ENABLED=0 GOOS=linux go build -o screenshot-service main.go

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
#ENV SERVER_HOST=http://192.168.1.254:1388  # Set your server IP here
CMD ["./screenshot-service"]