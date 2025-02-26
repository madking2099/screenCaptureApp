#!/usr/bin/env python3

import os
from fastapi import FastAPI, HTTPException
from fastapi.responses import RedirectResponse, JSONResponse
from fastapi.staticfiles import StaticFiles
from pydantic import BaseModel, HttpUrl
from typing import Dict, Optional
import uvicorn
from urllib.parse import urlparse
import uuid
from playwright.async_api import async_playwright
import asyncio
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

async def capture_webpage_screenshot(url: str, output_file: str, headers: Optional[Dict[str, str]] = None) -> str:
    try:
        logger.info(f"Fetching and screenshotting URL: {url}")
        async with async_playwright() as p:
            browser = await p.chromium.launch(headless=True)
            page = await browser.new_page()
            if headers:
                await page.set_extra_http_headers(headers)
            await page.goto(url, wait_until="networkidle", timeout=30000)  # 30s timeout
            await page.screenshot(path=output_file, full_page=True, type="png")
            await browser.close()
        logger.info(f"Screenshot saved to {output_file}")
        return output_file
    except Exception as e:
        logger.error(f"Screenshot failed: {str(e)}")
        raise HTTPException(status_code=500, detail=f"Error capturing screenshot: {str(e)}")

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
        screenshot_path = await capture_webpage_screenshot(str(request.url), output_file, request.headers)
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
