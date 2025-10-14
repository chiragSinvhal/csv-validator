# CSV Validator Service

> Professional Go REST API for CSV file processing with email validation

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shie## 📚 Documentation

| Document | Description |
|----------|-------------|
| **[Documentation Index](docs/DOCUMENTATION_INDEX.md)** | Complete documentation overview and navigation |
| **[API Reference](docs/API_REFERENCE.md)** | Complete API documentation with examples |
| **[Technical Overview](docs/TECHNICAL_OVERVIEW.md)** | Architecture and implementation details |
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](Dockerfile)

## 🎯 Purpose

A high-performance REST API service that processes CSV files to detect and flag rows containing valid email addresses.

## ✨ Key Features

- **📤 File Upload**: Secure CSV file upload with comprehensive validation
- **📧 Email Detection**: Intelligent email validation using regex patterns
- **⚡ Async Processing**: Non-blocking file processing with job status tracking
- **📥 File Download**: Retrieve processed files with proper HTTP status codes
- **🔒 Security**: Input validation, file type checking, and size limits
- **📊 Monitoring**: Health checks and structured logging
- **⚙️ Configurable Storage**: Environment-driven upload/download directory configuration
- **🐳 Docker Ready**: Containerized deployment with Docker Compose

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Git

### Installation & Setup
```bash
# Clone the repository
git clone <repository-url>
cd csv-validator

# Setup environment and dependencies
make setup

# Run the service
make run
```

The API will be available at `http://localhost:8080`

### Verify Installation
```bash
# Check health endpoint
curl http://localhost:8080/health

# Run integration tests
make integration-test
```

## 📋 Requirements

- Go 1.21 or higher
- Git

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

## 🎯 API Endpoints

### Upload CSV File
```http
POST /api/upload
Content-Type: multipart/form-data
```

**Request:** Form field `file` (CSV file)

**Success Response (200):**
```json
{
  "id": "a225eb00-0907-4273-92ca-5faadeefae5f"
}
```

**Error Response (400):**
```json
{
  "error": "Invalid file format. Only CSV files are allowed"
}
```

### Download Processed File
```http
GET /api/download/{job-id}
```

**Response Scenarios:**
- `200 OK`: Processed CSV file ready for download
- `423 Locked`: Job still processing
- `400 Bad Request`: Invalid job ID
- `404 Not Found`: File not found

### Health Check
```http
GET /health
```

Returns service health status.

## 🧪 Testing

```bash
# Run unit tests
make test

# Run tests with coverage
make coverage

# Run integration tests
make integration-test

# Run all checks (format, lint, test)
make check
```

## 🐳 Docker Deployment

### Quick Docker Run
```bash
# Build and run with Docker
make docker-build
make docker-run
```

### Docker Compose (Recommended)
```bash
# Run with Docker Compose (includes nginx proxy)
docker-compose up
```

### Manual Docker Commands
```bash
# Build image
docker build -t csv-validator .

# Run container
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e UPLOAD_DIR=./uploads \
  -e DOWNLOAD_DIR=./downloads \
  -e MAX_FILE_SIZE=10485760 \
  csv-validator
```

## ⚙️ Configuration

Configure via environment variables or `.env` file:

```env
PORT=8080                    # Server port
UPLOAD_DIR=./uploads         # Upload file storage directory
DOWNLOAD_DIR=./downloads     # Processed file storage directory
MAX_FILE_SIZE=10485760      # Max file size (10MB)
LOG_LEVEL=info              # Logging level (debug, info, warn, error)
GIN_MODE=release            # Framework mode (debug, release)
```

## 🔧 Development

### Development Commands
```bash
# Setup development environment
make setup

# Format code
make fmt

# Run linters
make lint

# Run go vet
make vet

# Install development tools
make install-tools
```

### Email Validation
The service validates email addresses using regex pattern:
```regex
^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
```

### CSV Processing Example

**Input CSV:**
```csv
name,email,age
Chirag,Chirag@example.com,30
Yash,invalid-email,25
Rohan,Rohan@test.org,35
```

**Output CSV:**
```csv
name,email,age,has_email
Chirag,Chirag@example.com,30,true
Yash,invalid-email,25,false
Rohan,Rohan@test.org,35,true
```

## � Documentation

| Document | Description |
|----------|-------------|
| **[API Reference](docs/API_REFERENCE.md)** | Complete API documentation with examples |
| **[Technical Overview](docs/TECHNICAL_OVERVIEW.md)** | Architecture and implementation details |

## 🛡️ Security Features

- **File Validation**: CSV format and size limits
- **Input Sanitization**: Filename and content validation
- **Path Protection**: Prevention of path traversal attacks
- **Content Validation**: MIME type and text file verification

## �📈 Performance Features

- **Async Processing**: Non-blocking file operations
- **Memory Efficient**: Streaming CSV parsing
- **Concurrent Safe**: Thread-safe job management
- **Resource Limits**: Configurable size and timeout limits

## � Monitoring & Logging

- **Structured Logging**: JSON format with configurable levels
- **Health Checks**: Service status monitoring
- **Request Tracking**: Request/response logging
- **Error Tracking**: Comprehensive error handling

## 🚨 Error Handling

### HTTP Status Codes
- `200 OK`: Successful request
- `400 Bad Request`: Invalid input (file format, job ID)
- `413 Payload Too Large`: File size exceeds limit
- `423 Locked`: Job still in progress
- `500 Internal Server Error`: Server processing error

### Common Errors
- Invalid file format (only CSV accepted)
- File size exceeds 10MB limit
- Malformed job ID format
- Empty or corrupted files
