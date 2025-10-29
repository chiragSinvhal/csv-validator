# Setup

## Requirements
Go 1.21+

## Quick Start

1. Clone and setup:
```bash
git clone <repo>
cd csv-validator
go mod download
```

2. Run:
```bash
go run .
```

3. Test:
```bash
curl http://localhost:8080/health
curl -X POST -F "file=@sample-data/sample1.csv" http://localhost:8080/api/upload
```

## Configuration

Copy `.env.example` to `.env` and edit as needed.

## Docker

```bash
docker-compose up
```
