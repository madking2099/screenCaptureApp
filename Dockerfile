FROM python:3.9-slim

WORKDIR /app

# Install dependencies and newer wkhtmltopdf/wkhtmltoimage
RUN apt-get update && apt-get install -y \
    wget \
    libxrender1 \
    libfontconfig1 \
    libxext6 \
    fonts-dejavu \
    && wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-3/wkhtmltox_0.12.6.1-3.buster_amd64.deb \
    && dpkg -i wkhtmltox_0.12.6.1-3.buster_amd64.deb || apt-get install -f -y \
    && rm wkhtmltox_0.12.6.1-3.buster_amd64.deb \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY screenshot_service.py .
RUN mkdir -p static

EXPOSE 8000

CMD ["uvicorn", "screenshot_service:app", "--host", "0.0.0.0", "--port", "8000"]
