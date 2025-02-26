# Webpage Screenshot Service

![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)
![Gin](https://img.shields.io/badge/Gin-1.9.1-green.svg)
![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)
![License](https://img.shields.io/badge/License-MIT-yellow.svg)

A lightweight microservice built with Go, Gin, and Docker to capture webpage screenshots as PNG images.  
It uses `chromedp` with headless Chromium for rendering, supporting modern web content.

## Features
- Gin Backend: Fast, minimalist Go web framework.
- Dockerized: Ready-to-run container with all dependencies.
- Custom Headers: Supports optional HTTP headers for requests.
- Robust Rendering: Uses `chromedp` for full JS/CSS support.
- API Endpoints:
  - `POST /screenshot`: Capture a webpage screenshot.
  - `GET /health`: Check service status.
  - `DELETE /static/<filename>`: Remove a screenshot file.

## Prerequisites
- Docker installed on your system.
- Basic knowledge of command-line tools and HTTP requests.

## Installation

### Clone the Repository
Run: git clone https://github.com/madking2099/screenCaptureApp.git  
Then: cd screenCaptureApp

### Build the Docker Image
Run: docker build -t screenshot-service:latest .

### Run the Container
Run: docker run -d -p 1388:8000 --name screenshot-service screenshot-service:latest  
The service will be available at `http://localhost:1388`.

## Usage

### Health Check
Run: curl http://localhost:1388/health  
Response: {"status": "healthy"}

### Capture a Screenshot

#### Via curl
Run: curl -X POST "http://localhost:1388/screenshot" -H "Content-Type: application/json" -d '{"url": "https://example.com"}'  
Response: {"file_url": "/static/screenshot_<id>.png"}  
Download: curl "http://localhost:1388<file_url>" --output screenshot.png

#### With Basic Auth
Run: curl -X POST "http://localhost:1388/screenshot" -H "Content-Type: application/json" -d '{"url": "https://user:pass@example.com"}'  
Response: {"file_url": "/static/screenshot_<id>.png"}

### Delete a Screenshot
Run: curl -X DELETE "http://localhost:1388/static/<filename>"  
Response: {"message": "File <filename> deleted"}

## API Reference

### POST /screenshot
Request Body:  
{  
  "url": "https://example.com",  
  "headers": {"User-Agent": "MyBot/1.0"},  
  "output_filename": "screenshot"  
}  
Response: {"file_url": "/static/screenshot.png"}

### GET /health
Response: {"status": "healthy"}

### DELETE /static/<filename>
Response: {"message": "File <filename> deleted"}

## Development

### Dependencies
- Go 1.21+
- `chromedp` (Chromium driver)
- `gin` (web framework)

### Local Development
Run: go mod tidy  
Then: go run main.go

## Limitations
- Self-signed certs are bypassed (insecure); use a CA cert for production.
- Image size is ~250-300MB due to Chromium.

## License
This project is licensed under the MIT License.  
See the `LICENSE` file for details.

## Contributing
Feel free to open issues or submit pull requests! Contributions are welcome.

## Acknowledgments
- Built with `Gin` (https://gin-gonic.com/).
- Powered by `chromedp` (https://github.com/chromedp/chromedp).
- Inspired by a need for reliable screenshot tools.

---
Developed by Grok 3 (xAI) · February 2025