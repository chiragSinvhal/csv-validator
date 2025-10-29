# Getting Started Guide

Quick setup instructions for the CSV Validator Service.

## Prerequisites

- Go 1.21 or later installed
- Basic familiarity with REST APIs

## Setup

1. **Clone the repository**
```bash
git clone <repository-url>
cd csv-validator
```

2. **Install dependencies**
```bash
go mod download
```

3. **Run the service**
```bash
go run .
```

The service will start on port 8080.

## First Test

1. **Check if the service is running**
```bash
curl http://localhost:8080/health
```

Expected response: `{"status":"healthy"}`

2. **Upload a test file**
```bash
curl -X POST -F "file=@sample-data/sample1.csv" http://localhost:8080/api/upload
```

You'll get a response with a job ID.

3. **Download the processed file**
```bash
curl http://localhost:8080/api/download/YOUR_JOB_ID -o processed.csv
```

## What Happens

The service processes your CSV file and adds a new column called `has_email` that indicates whether each row contains a valid email address.

## Docker Alternative

If you prefer using Docker:

```bash
docker-compose up
```

This will start the service with all dependencies configured.

## Next Steps

- Check the API documentation for detailed endpoint information
- Look at the sample CSV files in `sample-data/` for examples
- Review the configuration options in `.env.example`
