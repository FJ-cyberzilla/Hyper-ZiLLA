# Use official Node.js base image
FROM node:18-slim

# Set environment variables
ENV NODE_ENV=production
ENV PYTHONUNBUFFERED=1

# Set working directory
WORKDIR /usr/src/app

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 \
    python3-pip \
    python3-venv \
    postgresql-client \
    libopencv-dev \
    tesseract-ocr \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies
COPY requirements.txt .
RUN pip3 install --upgrade pip && pip3 install --no-cache-dir -r requirements.txt && rm requirements.txt

# Copy Node.js package files and install dependencies
COPY package*.json ./
RUN npm ci --omit=dev

# Copy application code
COPY . .

# Create necessary directories
RUN mkdir -p logs data temp

# Security hardening
RUN chmod -R 750 /usr/src/app \
    && useradd -r -s /bin/false zilla \
    && chown -R zilla:zilla /usr/src/app

USER zilla

# Expose ports
EXPOSE 3000 9229

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD node healthcheck.js || exit 1

# Start command
CMD ["node", "--max-old-space-size=4096", "main.js"]
