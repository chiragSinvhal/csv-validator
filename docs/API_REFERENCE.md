# CSV Validator API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
No authentication required for this API.

## Endpoints

### Health Check

#### GET /health
Check if the service is running.

**Response:**
```json
{
  "status": "healthy"
}
```

### File Upload

#### POST /api/upload
Upload a CSV file for email validation processing.

**Request:**
- Method: `POST`
- Content-Type: `multipart/form-data`
- Body: Form data with file field named `file`

**Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| file | File | Yes | CSV file to process (max 10MB) |

**Response (Success - 200 OK):**
```json
{
  "id": "a225eb00-0907-4273-92ca-5faadeefae5f"
}
```

**Response (Error - 400 Bad Request):**
```json
{
  "error": "Invalid file format. Only CSV files are allowed"
}
```

**Response (Error - 413 Payload Too Large):**
```json
{
  "error": "File size (15728640 bytes) exceeds maximum allowed size (10485760 bytes)"
}
```

**Response (Error - 500 Internal Server Error):**
```json
{
  "error": "Failed to save uploaded file"
}
```

**Example using curl:**
```bash
curl -X POST \
  http://localhost:8080/api/upload \
  -H 'Content-Type: multipart/form-data' \
  -F 'file=@sample.csv'
```

### File Download

#### GET /api/download/{id}
Download the processed CSV file.

**Request:**
- Method: `GET`
- Path Parameter: `id` (UUID of the processing job)

**Response (Success - 200 OK):**
- Content-Type: `text/csv`
- Content-Disposition: `attachment; filename=processed_filename.csv`
- Body: Processed CSV file with email validation flags

**Response (Job In Progress - 423 Locked):**
```json
{
  "error": "Job is still in progress"
}
```

**Response (Invalid ID - 400 Bad Request):**
```json
{
  "error": "Invalid job ID format"
}
```

**Response (Job Not Found - 400 Bad Request):**
```json
{
  "error": "Invalid job ID"
}
```

**Response (Job Failed - 500 Internal Server Error):**
```json
{
  "error": "Job processing failed: {failure reason}"
}
```

**Response (File Not Found - 404 Not Found):**
```json
{
  "error": "Processed file not found"
}
```

**Example using curl:**
```bash
curl -X GET \
  http://localhost:8080/api/download/a225eb00-0907-4273-92ca-5faadeefae5f \
  -o processed_file.csv
```

## CSV Processing Details

### Email Validation
The service validates email addresses using the following regex pattern:
```
^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
```

### Output Format
For each row in the input CSV (excluding the header), a new column `has_email` is added:
- `true`: The row contains at least one valid email address
- `false`: The row does not contain any valid email addresses

### Example

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

## Error Handling

### HTTP Status Codes
- `200 OK`: Successful request
- `400 Bad Request`: Invalid request (bad file format, invalid job ID, etc.)
- `413 Payload Too Large`: File size exceeds maximum limit
- `423 Locked`: Job is still in progress (for download requests)
- `404 Not Found`: Processed file not found
- `500 Internal Server Error`: Server error during processing

### File Validation
- Only `.csv` files are accepted
- Maximum file size: 10MB (configurable)
- Files must contain valid text content
- Empty files are rejected

## Rate Limiting
Currently, no rate limiting is implemented. In production, consider implementing rate limiting based on your requirements.

## CORS
CORS is enabled for all origins (`*`). In production, configure appropriate CORS policies.

## Examples

### Complete Workflow

1. **Upload a file:**
```bash
curl -X POST \
  http://localhost:8080/api/upload \
  -H 'Content-Type: multipart/form-data' \
  -F 'file=@sample.csv'
```

Response:
```json
{
  "id": "a225eb00-0907-4273-92ca-5faadeefae5f"
}
```

2. **Check if processing is complete:**
```bash
curl -X GET \
  http://localhost:8080/api/download/a225eb00-0907-4273-92ca-5faadeefae5f
```

If still processing (423 response):
```json
{
  "error": "Job is still in progress"
}
```

3. **Download processed file when ready:**
```bash
curl -X GET \
  http://localhost:8080/api/download/a225eb00-0907-4273-92ca-5faadeefae5f \
  -o processed_sample.csv
```

### Using JavaScript/Fetch API

```javascript
// Upload file
const formData = new FormData();
formData.append('file', fileInput.files[0]);

fetch('http://localhost:8080/api/upload', {
    method: 'POST',
    body: formData
})
.then(response => response.json())
.then(data => {
    const jobId = data.id;
    // Poll for completion
    checkJobStatus(jobId);
});

// Check job status and download when ready
function checkJobStatus(jobId) {
    fetch(`http://localhost:8080/api/download/${jobId}`)
    .then(response => {
        if (response.status === 200) {
            // File ready, download it
            return response.blob();
        } else if (response.status === 423) {
            // Still processing, check again later
            setTimeout(() => checkJobStatus(jobId), 1000);
        } else {
            throw new Error('Job failed or invalid');
        }
    })
    .then(blob => {
        if (blob) {
            // Create download link
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'processed_file.csv';
            a.click();
        }
    });
}
```