# CSV Validator Service

A REST API service for CSV file processing with email validation functionality.

## Overview

This service processes CSV files to detect and flag rows containing valid email addresses. It provides a simple HTTP API for uploading CSV files and downloading processed results.

## Features

- CSV file upload with validation
- Email detection and flagging
- Asynchronous file processing
- Job status tracking
- File download with proper HTTP status codes
- Input validation and security checks
- Health monitoring endpoint
- Docker support

## 🏗️ Project Structure

```
csv-validator/
├── cmd/server/              # Application entry point
├── internal/                # Private application code
│   ├── config/             # Configuration management
│   ├── handlers/           # HTTP request handlers
│   ├── models/             # Data models
│   ├── services/           # Business logic
│   └── utils/              # Utility functions
├── pkg/logger/             # Logging package
├── docs/                   # Documentation files
├── scripts/                # Build and test scripts
├── sample-data/            # Sample CSV files for testing
└── Docker files & configs
```

## Quick Start

### Requirements
- Go 1.21 or higher

### Installation
```bash
git clone <repository-url>
cd csv-validator
go mod download
```

### Running the Service
```bash
go run .
```

The service runs on `http://localhost:8080` by default.

### Basic Usage
```bash
# Check if service is running
curl http://localhost:8080/health

# Upload a CSV file
curl -X POST -F "file=@sample.csv" http://localhost:8080/api/upload

# Download processed file (replace {job-id} with actual ID from upload response)
curl http://localhost:8080/api/download/{job-id} -o processed.csv
```

## API Documentation

### Upload Endpoint
**POST /api/upload**

Upload a CSV file for processing.

Request:
- Content-Type: multipart/form-data
- Field: `file` (CSV file, max 10MB)

Response:
```json
{
  "id": "uuid-job-id"
}
```

### Download Endpoint
**GET /api/download/{id}**

Download the processed CSV file.

Responses:
- 200: File ready for download
- 423: Processing still in progress
- 400: Invalid job ID
- 404: File not found

### Health Check
**GET /health**

Check service status.

## Configuration

Set environment variables or create a `.env` file:

```env
PORT=8080
UPLOAD_DIR=./uploads
DOWNLOAD_DIR=./downloads
MAX_FILE_SIZE=10485760
LOG_LEVEL=info
GIN_MODE=release
```

## Docker

```bash
# Build and run
docker build -t csv-validator .
docker run -p 8080:8080 csv-validator

# Or use docker-compose
docker-compose up
```

## Development

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Format code
go fmt ./...

# Build
go build -o csv-validator .
```

## How It Works

The service adds a `has_email` column to uploaded CSV files:
- `true` if the row contains at least one valid email address
- `false` if no valid email addresses are found

Example:
```csv
# Input
name,email,phone
John,john@example.com,123-456-7890
Jane,,987-654-3210

# Output
name,email,phone,has_email
John,john@example.com,123-456-7890,true
Jane,,987-654-3210,false
```

## Error Handling

Common error responses:
- Invalid file format (non-CSV files)
- File size exceeds 10MB limit
- Invalid job ID format
- Processing failures

## Testing

Sample CSV files are provided in the `sample-data/` directory for testing.

For automated testing, run the integration test script:
```bash
./scripts/integration-tests.sh
```
