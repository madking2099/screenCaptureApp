#!/usr/bin/env python3

import requests
import subprocess
import os
from fastapi import FastAPI, HTTPException, Response
from fastapi.responses import RedirectResponse  # Add this import
from pydantic import BaseModel, HttpUrl
from typing import Dict, Optional
import uvicorn
from urllib.parse import urlparse

app = FastAPI(
    title="Webpage Screenshot Service",
    description="Capture screenshots of webpages as images.",
    version="1.0.0"
)

class ScreenshotRequest(BaseModel):
    url: HttpUrl
    headers: Optional[Dict[str, str]] = None
    output_filename: Optional[str] = "screenshot.png"

def is_valid_url(url: str) -> bool:
    try:
        result = urlparse(url)
        return all([result.scheme, result.netloc])
    except ValueError:
        return False

def capture_webpage_screenshot(url: str, output_file: str = "screenshot.png", headers: Optional[Dict[str, str]] = None) -> str:
    temp_html = "temp.html"
    default_headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    if headers:
        default_headers.update(headers)

    try:
        response = requests.get(url, headers=default_headers, timeout=10)
        response.raise_for_status()
        with open(temp_html, "w", encoding="utf-8") as f:
            f.write(response.text)
        subprocess.run([
            "wkhtmltoimage",
            "--width", "1280",
            "--quality", "90",
            temp_html,
            output_file
        ], check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        return output_file
    except requests.RequestException as e:
        raise HTTPException(status_code=400, detail=f"Error fetching webpage: {str(e)}")
    except subprocess.CalledProcessError as e:
        raise HTTPException(status_code=500, detail=f"Error rendering screenshot: {e.stderr.decode().strip()}")
    finally:
        if os.path.exists(temp_html):
            os.remove(temp_html)

# New root route to redirect to Swagger UI
@app.get("/", response_class=RedirectResponse)
async def root():
    return "/docs"

@app.post("/screenshot/", response_class=Response)
async def create_screenshot(request: ScreenshotRequest):
    if not is_valid_url(str(request.url)):
        raise HTTPException(status_code=400, detail="Invalid URL provided.")
    output_file = request.output_filename if request.output_filename.endswith(".png") else f"{request.output_filename}.png"
    try:
        screenshot_path = capture_webpage_screenshot(str(request.url), output_file, request.headers)
        with open(screenshot_path, "rb") as f:
            image_data = f.read()
        if os.path.exists(screenshot_path):
            os.remove(screenshot_path)
        return Response(content=image_data, media_type="image/png")
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Server error: {str(e)}")

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
