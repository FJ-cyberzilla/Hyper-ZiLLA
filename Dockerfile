FROM python:3.10-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Create necessary directories
RUN mkdir -p logs data cache

# Expose port
EXPOSE 5000

# Set environment variables
ENV PYTHONPATH=/app
ENV HZ_WEB_HOST=0.0.0.0
ENV HZ_WEB_PORT=5000

# Run application
CMD ["python", "main.py", "--mode", "web"]
