FROM python:3.9-slim

WORKDIR /app

RUN apt-get update && apt-get install -y \
    wkhtmltopdf \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY screenshot_service.py .

# Create static directory (optional, as the app creates it, but ensures it exists)
RUN mkdir -p static

EXPOSE 8000

CMD ["uvicorn", "screenshot_service:app", "--host", "0.0.0.0", "--port", "8000"]
