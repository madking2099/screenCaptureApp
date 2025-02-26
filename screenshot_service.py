#!/usr/bin/env python3

import requests
import subprocess
import os
from fastapi import FastAPI, HTTPException
from fastapi.responses import RedirectResponse, JSONResponse
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel, HttpUrl
from typing import Dict, Optional
import uvicorn
from urllib.parse import urlparse
import uuid
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Webpage Screenshot Service",
    description="Capture screenshots of webpages as images and serve them.",
    version="1.0.0"
)

app.mount("/static", StaticFiles(directory="static"), name="static")

class ScreenshotRequest(BaseModel):
    url: HttpUrl
    headers: Optional[Dict[str, str]] = None
    output_filename: Optional[str] = None

def is_valid_url(url: str) -> bool:
    try:
        result = urlparse(url)
        return all([result.scheme, result.netloc])
    except ValueError:
        return False

def capture_webpage_screenshot(url: str, output_file: str, headers: Optional[Dict[str, str]] = None) -> str:
    temp_html = "temp.html"
    default_headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    }
    if headers:
        default_headers.update(headers)

    try:
        logger.info(f"Fetching URL: {url}")
        response = requests.get(url, headers=default_headers, timeout=10)
        response.raise_for_status()
        with open(temp_html, "w", encoding="utf-8") as f:
            f.write(response.text)
        logger.info(f"Rendering {temp_html} to {output_file}")
        result = subprocess.run([
            "wkhtmltoimage",
            "--width", "1280",
            "--quality", "90",
            "--enable-javascript",        # Enable JS rendering
            "--javascript-delay", "2000", # Wait for JS
            "--no-stop-slow-scripts",     # Don’t halt on slow JS
            "--ignore-load-errors",       # Skip protocol errors
            temp_html,
            output_file
        ], check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        logger.info(f"wkhtmltoimage output: {result.stdout.decode()}")
        if result.stderr:
            logger.warning(f"wkhtmltoimage warnings: {result.stderr.decode()}")
        return output_file
    except requests.RequestException as e:
        logger.error(f"Request failed: {str(e)}")
        raise HTTPException(status_code=400, detail=f"Error fetching webpage: {str(e)}")
    except subprocess.CalledProcessError as e:
        logger.error(f"wkhtmltoimage failed: {e.stderr.decode().strip()}")
        raise HTTPException(status_code=500, detail=f"Error rendering screenshot: {e.stderr.decode().strip()}")
    finally:
        if os.path.exists(temp_html):
            os.remove(temp_html)

@app.get("/", response_class=RedirectResponse)
async def root():
    return "/docs"

@app.post("/screenshot/")
async def create_screenshot(request: ScreenshotRequest):
    if not is_valid_url(str(request.url)):
        raise HTTPException(status_code=400, detail="Invalid URL provided.")
    
    if not os.path.exists("static"):
        os.makedirs("static")
    
    filename = request.output_filename if request.output_filename else f"screenshot_{uuid.uuid4().hex}"
    output_file = f"static/{filename}.png" if not filename.endswith(".png") else f"static/{filename}"
    
    try:
        screenshot_path = capture_webpage_screenshot(str(request.url), output_file, request.headers)
        file_url = f"/static/{os.path.basename(screenshot_path)}"
        return JSONResponse(content={"file_url": file_url}, status_code=200)
    except Exception as e:
        if os.path.exists(output_file):
            os.remove(output_file)
        raise

@app.delete("/static/{filename}")
async def delete_screenshot(filename: str):
    file_path = f"static/{filename}"
    if not os.path.exists(file_path):
        raise HTTPException(status_code=404, detail="File not found")
    try:
        os.remove(file_path)
        return JSONResponse(content={"message": f"File {filename} deleted"}, status_code=200)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error deleting file: {str(e)}")

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
