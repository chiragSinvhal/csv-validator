# Getting Started

> Quick setup guide to run the CSV Validator Service

## 🚀 Quick Start

### Prerequisites
- Go 1.21+ installed
- Git installed

### 1. Clone and Setup
```bash
git clone https://github.com/chiragSinvhal/csv-validator.git
cd csv-validator
go mod download
```

### 2. Run the Service
```bash
go run cmd/server/main.go
```

The service will start on `http://localhost:8080`

### 3. Test the Service

**Check if it's running:**
```bash
curl http://localhost:8080/health
```
Expected response: `{"status":"OK","timestamp":"..."}`

**Upload a CSV file:**
```bash
curl -X POST -F "file=@sample-data/sample1.csv" http://localhost:8080/api/upload
```
Expected response: `{"id":"some-uuid-here"}`

**Download processed file:**
```bash
curl http://localhost:8080/api/download/{your-job-id} -o processed.csv
```

### 4. Expected Output

The service adds an `has_email` column to your CSV:
- `true` - Row contains at least one valid email
- `false` - Row contains no valid emails

**Example:**
```csv
name,email,phone,has_email
John,john@email.com,123-456-7890,true
Jane,,987-654-3210,false
```

## 🐳 Alternative: Using Docker

```bash
docker-compose up
```

## 📚 More Information

- **[Complete Documentation](docs/DOCUMENTATION_INDEX.md)** - Full project details
- **[API Reference](docs/API_REFERENCE.md)** - Detailed API documentation
- **[README.md](README.md)** - Complete project overview
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
  -e MAX_FILE_SIZE=10485760 \
  csv-validator
```

## ⚙️ Configuration

Configure via environment variables or `.env` file:

```env
PORT=8080                    # Server port
UPLOAD_DIR=./uploads         # File storage directory
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
