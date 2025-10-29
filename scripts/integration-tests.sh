#!/bin/bash

# CSV Validator API Integration Tests

API_BASE="http://localhost:8080"
TEST_FILE="sample-data/sample1.csv"

echo "CSV Validator API Integration Tests"
echo "==================================="
echo ""

# Check if server is running
echo "Testing health endpoint..."
health_response=$(curl -s -w "%{http_code}" -o /dev/null "$API_BASE/health")
if [ "$health_response" = "200" ]; then
    echo "Health check: PASS"
else
    echo "Health check: FAIL (HTTP $health_response)"
    echo "Make sure the server is running: go run ."
    exit 1
fi
echo ""

# Test file upload
echo "Testing file upload..."
if [ ! -f "$TEST_FILE" ]; then
    echo "Test file not found: $TEST_FILE"
    exit 1
fi

upload_response=$(curl -s -X POST \
    -F "file=@$TEST_FILE" \
    "$API_BASE/api/upload")

job_id=$(echo "$upload_response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

if [ -n "$job_id" ]; then
    echo "File upload: PASS"
    echo "Job ID: $job_id"
else
    echo "File upload: FAIL"
    echo "Response: $upload_response"
    exit 1
fi
echo ""

# Test job status and download
echo "Testing file processing..."
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    download_response=$(curl -s -w "%{http_code}" \
        "$API_BASE/api/download/$job_id" \
        -o "test_output.csv")
    
    http_code=${download_response: -3}
    
    if [ "$http_code" = "200" ]; then
        echo "File processing: PASS"
        echo "Output saved as: test_output.csv"
        break
    elif [ "$http_code" = "423" ]; then
        echo "Processing... (attempt $((attempt + 1))/$max_attempts)"
        sleep 1
        attempt=$((attempt + 1))
    else
        echo "File processing: FAIL (HTTP $http_code)"
        exit 1
    fi
done

if [ $attempt -eq $max_attempts ]; then
    echo "File processing: TIMEOUT"
    exit 1
fi
echo ""

# Verify output file
echo "Verifying output file..."
if [ -f "test_output.csv" ]; then
    if grep -q "has_email" "test_output.csv"; then
        echo "Output verification: PASS"
        echo ""
        echo "Sample output:"
        head -3 "test_output.csv"
    else
        echo "Output verification: FAIL (missing has_email column)"
        exit 1
    fi
else
    echo "Output verification: FAIL (file not found)"
    exit 1
fi
echo ""

# Test error cases
echo "Testing error handling..."

# Invalid file type
echo "This is not a CSV" > test_invalid.txt
invalid_response=$(curl -s -X POST \
    -F "file=@test_invalid.txt" \
    "$API_BASE/api/upload")

if echo "$invalid_response" | grep -q "error"; then
    echo "Error handling: PASS"
else
    echo "Error handling: FAIL"
fi
rm -f test_invalid.txt

echo ""
echo "All tests completed successfully!"
echo ""
echo "Cleaning up..."
rm -f test_output.csv
echo "Done."
