# Use official Python runtime as base image
FROM python:3.9-slim

# Set working directory
WORKDIR /app

# Install system dependencies (wkhtmltopdf for wkhtmltoimage)
RUN apt-get update && apt-get install -y \
    wkhtmltopdf \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements file
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the application code
COPY screenshot_service.py .

# Expose port 8000
EXPOSE 8000

# Run the FastAPI app with uvicorn
CMD ["uvicorn", "screenshot_service:app", "--host", "0.0.0.0", "--port", "8000"]