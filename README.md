# CSV Email Checker

Quick service I built to scan CSV files for email addresses. Pretty straightforward - upload a CSV, get it back with a new column showing which rows have emails.

## What it does

Takes your CSV file and adds a `has_email` column. That's it. If a row has what looks like an email address anywhere in it, it gets marked as `true`, otherwise `false`.

## Running it

You'll need Go 1.21+. Then just:

```bash
git clone <your-repo>
cd csv-validator
go run .
```

Runs on port 8080 by default.

## Usage

Upload a file:
```bash
curl -X POST -F "file=@your-file.csv" http://localhost:8080/api/upload
```

You'll get back a job ID. Use that to download:
```bash
curl http://localhost:8080/api/download/{job-id} -o result.csv
```

Note: If you try to download too quickly, you'll get a 423 status while it's still processing.

## Example

Input CSV:
```
name,contact,notes
Alice,alice@gmail.com,good customer
Bob,555-1234,called yesterday
```

Output:
```
name,contact,notes,has_email
Alice,alice@gmail.com,good customer,true
Bob,555-1234,called yesterday,false
```

## Config

Set these if you want to change defaults:
- `PORT` - server port (default: 8080)
- `MAX_FILE_SIZE` - max upload size in bytes (default: 10MB)
- `UPLOAD_DIR` - where to store uploads (default: ./uploads)

## Docker

```bash
docker build -t csv-validator .
docker run -p 8080:8080 csv-validator
```

## Known issues

- Files are processed synchronously right now, so big files might timeout
- Email regex could probably be better, but works for most common formats
- No cleanup of old processed files yet

## Testing

```bash
go test ./...
```

There's some sample data in `sample-data/` if you need test files.

## Health check

Hit `/health` to see if it's running.
