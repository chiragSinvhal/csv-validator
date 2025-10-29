# Technical Overview

## Project Summary

A Go-based REST API service for processing CSV files with email validation.

## Architecture

### Core Components

**Services Layer**
- FileService: Handles file storage and validation
- JobService: Manages job lifecycle and status
- CSVService: Processes CSV files and adds email validation

**HTTP Layer**
- Handlers: Process HTTP requests and responses
- Middleware: CORS, logging, error handling

**Configuration**
- Environment-based configuration
- File storage directory management

### Key Design Decisions

**Asynchronous Processing**
File processing happens in background goroutines to avoid blocking HTTP requests.

**In-Memory Job Storage**
Jobs are stored in memory with mutex protection for thread safety. This works well for single-instance deployments.

**File System Storage**
Uploaded and processed files are stored on the local filesystem with unique naming to prevent conflicts.

## Email Validation

Uses regex pattern matching to identify valid email addresses:
```regex
^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
```

## Security Measures

- File type validation (CSV only)
- File size limits (configurable, default 10MB)
- Input sanitization for filenames
- UUID-based job IDs to prevent enumeration
- Content validation to ensure text files

## Error Handling

The service uses proper HTTP status codes:
- 200: Success
- 400: Bad request (invalid input)
- 413: File too large
- 423: Job still processing
- 500: Server error

## Configuration

Environment variables control:
- Server port
- Upload/download directories
- File size limits
- Logging levels
- Framework mode

## Testing Strategy

- Unit tests for all services
- Integration tests for HTTP endpoints
- End-to-end API testing script
- File validation test cases

## Deployment

**Development**: Direct Go execution
**Production**: Docker containers with optional nginx proxy

## Performance Considerations

- Streaming CSV parsing for memory efficiency
- Concurrent job processing
- Configurable timeouts and limits
- Graceful shutdown handling

## Monitoring

- Health check endpoint
- Structured logging
- Request/response tracking
- Error tracking and reporting

## Future Enhancements

For production scale, consider:
- Database-backed job storage
- Distributed file storage
- Message queue for processing
- Horizontal scaling support
