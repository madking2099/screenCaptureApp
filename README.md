# Webpage Screenshot Service

![Python](https://img.shields.io/badge/Python-3.9+-blue.svg)
![FastAPI](https://img.shields.io/badge/FastAPI-0.95.1-green.svg)
![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)
![License](https://img.shields.io/badge/License-MIT-yellow.svg)

A lightweight microservice built with FastAPI and Docker to capture webpage screenshots as PNG images.  
It uses `wkhtmltoimage` for rendering, avoiding the overhead of full browser automation tools like Selenium.

## Features
- FastAPI Backend: Modern, async Python web framework.
- Dockerized: Ready-to-run container with all dependencies.
- Custom Headers: Supports optional HTTP headers for requests.
- Lightweight: Uses `wkhtmltoimage` instead of heavy browser engines.
- API Endpoints:
  - `POST /screenshot/`: Capture a webpage screenshot.
  - `GET /health`: Check service status.

## Prerequisites
- Docker installed on your system.
- Basic knowledge of command-line tools and HTTP requests.

## Installation

### Clone the Repository
Run: git clone https://github.com/madking2099/screenCaptureApp.git  
Then: cd screenshot-service

### Build the Docker Image
Run: docker build -t screenshot-service:latest .

### Run the Container
Run: docker run -d -p 8000:8000 --name screenshot-service screenshot-service:latest  
The service will be available at `http://localhost:8000`.

## Usage

### Health Check
Run: curl http://localhost:8000/health  
Response: {"status": "healthy"}

### Capture a Screenshot

#### Via curl
Run: curl -X POST "http://localhost:8000/screenshot/" -H "Content-Type: application/json" -d '{"url": "https://example.com"}' --output screenshot.png

#### Via Python
Example script:
import requests  
payload = {  
    "url": "https://example.com",  
    "headers": {"User-Agent": "MyBot/1.0"},  
    "output_filename": "example_shot"  
}  
response = requests.post("http://localhost:8000/screenshot/", json=payload)  
if response.status_code == 200:  
    with open("example_shot.png", "wb") as f:  
        f.write(response.content)  
    print("Screenshot saved!")  
else:  
    print(f"Error: {response.text}")

## API Reference

### POST /screenshot/
Request Body:  
{  
  "url": "https://example.com",  
  "headers": {"User-Agent": "MyBot/1.0"},  
  "output_filename": "screenshot"  
}  
Response: Raw PNG image (`Content-Type: image/png`).

### GET /health
Response: {"status": "healthy"}

## Development

### Dependencies
- Python 3.9+
- FastAPI, Uvicorn, Requests, Pydantic (see `requirements.txt`)
- `wkhtmltopdf` (for `wkhtmltoimage`)

### Local Development
Run: pip install -r requirements.txt  
Then: uvicorn screenshot_service:app --host 0.0.0.0 --port 8000 --reload

## Limitations
- JavaScript Rendering: `wkhtmltoimage` may not fully render dynamic, JS-heavy pages (e.g., SPAs).  
  Consider alternatives like Puppeteer for such cases.
- File Cleanup: Screenshots are deleted after serving; modify the code for persistence if needed.

## License
This project is licensed under the MIT License.  
See the `LICENSE` file for details.

## Contributing
Feel free to open issues or submit pull requests! Contributions are welcome.

## Acknowledgments
- Built with `FastAPI` (https://fastapi.tiangolo.com/).
- Powered by `wkhtmltopdf` (https://wkhtmltopdf.org/).
- Inspired by a need for lightweight screenshot tools.

---
Developed by Grok 3 (xAI) · February 2025
