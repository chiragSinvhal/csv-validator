# API Reference

Complete API documentation for the CSV Validator Service.

## Base URL
```
http://localhost:8080
```

## Endpoints

### Health Check

**GET /health**

Check service health status.

**Response:**
```json
{
  "status": "healthy"
}
```

### File Upload

**POST /api/upload**

Upload a CSV file for email validation processing.

**Request:**
- Method: POST
- Content-Type: multipart/form-data
- Body: Form data with file field named `file`

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| file | File | Yes | CSV file (max 10MB) |

**Success Response (200):**
```json
{
  "id": "a225eb00-0907-4273-92ca-5faadeefae5f"
}
```

**Error Responses:**

400 Bad Request:
```json
{
  "error": "Invalid file format. Only CSV files are allowed"
}
```

413 Payload Too Large:
```json
{
  "error": "File size exceeds maximum allowed size"
}
```

**Example:**
```bash
curl -X POST \
  http://localhost:8080/api/upload \
  -F 'file=@sample.csv'
```

### File Download

**GET /api/download/{id}**

Download the processed CSV file.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| id | string | Yes | Job ID from upload response |

**Success Response (200):**
- Content-Type: text/csv
- Body: Processed CSV file with has_email column

**Error Responses:**

423 Locked (Processing):
```json
{
  "error": "Job is still in progress"
}
```

400 Bad Request:
```json
{
  "error": "Invalid job ID format"
}
```

404 Not Found:
```json
{
  "error": "Processed file not found"
}
```

**Example:**
```bash
curl -X GET \
  http://localhost:8080/api/download/a225eb00-0907-4273-92ca-5faadeefae5f \
  -o processed_file.csv
```

## Processing Details

### Email Validation

The service validates emails using this pattern:
```
^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$
```

### Output Format

A new column `has_email` is added to each row:
- `true`: Row contains at least one valid email
- `false`: No valid emails found

### Example

Input:
```csv
name,email,age
John,john@example.com,30
Jane,invalid-email,25
```

Output:
```csv
name,email,age,has_email
John,john@example.com,30,true
Jane,invalid-email,25,false
```

## Error Handling

### Status Codes
- 200: Success
- 400: Bad request
- 413: File too large
- 423: Processing in progress
- 404: Not found
- 500: Server error

### File Validation
- Only .csv files accepted
- Maximum 10MB file size
- Must contain valid text content
- Empty files rejected

## Complete Workflow Example

1. Upload file:
```bash
curl -X POST -F 'file=@data.csv' http://localhost:8080/api/upload
```

Response:
```json
{"id": "job-uuid-here"}
```

2. Check status:
```bash
curl http://localhost:8080/api/download/job-uuid-here
```

If still processing, you'll get HTTP 423.

3. Download when ready:
```bash
curl http://localhost:8080/api/download/job-uuid-here -o result.csv
```

## JavaScript Example

```javascript
async function processCSV(file) {
    // Upload
    const formData = new FormData();
    formData.append('file', file);
    
    const uploadResponse = await fetch('/api/upload', {
        method: 'POST',
        body: formData
    });
    
    const { id } = await uploadResponse.json();
    
    // Poll for completion
    while (true) {
        const downloadResponse = await fetch(`/api/download/${id}`);
        
        if (downloadResponse.status === 200) {
            // Ready - download file
            const blob = await downloadResponse.blob();
            // Handle download
            break;
        } else if (downloadResponse.status === 423) {
            // Still processing - wait and retry
            await new Promise(resolve => setTimeout(resolve, 1000));
        } else {
            // Error
            throw new Error('Processing failed');
        }
    }
}
```
