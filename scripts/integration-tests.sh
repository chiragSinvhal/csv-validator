#!/bin/bash

echo "Testing CSV Validator..."

API="http://localhost:8080"
TEST_FILE="sample-data/sample1.csv"

# Health check
echo "Checking health..."
if curl -s $API/health | grep -q "healthy"; then
    echo "Health: OK"
else
    echo "Health: FAIL"
    exit 1
fi

# Upload test
echo "Testing upload..."
if [ ! -f "$TEST_FILE" ]; then
    echo "Test file missing: $TEST_FILE"
    exit 1
fi

response=$(curl -s -X POST -F "file=@$TEST_FILE" $API/api/upload)
job_id=$(echo $response | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

if [ -n "$job_id" ]; then
    echo "Upload: OK (job: $job_id)"
else
    echo "Upload: FAIL"
    echo "Response: $response"
    exit 1
fi

# Download test
echo "Testing download..."
attempts=0
max_attempts=20

while [ $attempts -lt $max_attempts ]; do
    response=$(curl -s -w "%{http_code}" $API/api/download/$job_id -o test_output.csv)
    code=${response: -3}
    
    if [ "$code" = "200" ]; then
        echo "Download: OK"
        break
    elif [ "$code" = "423" ]; then
        echo "Processing... ($((attempts + 1))/$max_attempts)"
        sleep 1
        attempts=$((attempts + 1))
    else
        echo "Download: FAIL (HTTP $code)"
        exit 1
    fi
done

if [ $attempts -eq $max_attempts ]; then
    echo "Download: TIMEOUT"
    exit 1
fi

# Verify output
if [ -f "test_output.csv" ] && grep -q "has_email" "test_output.csv"; then
    echo "Verification: OK"
    echo "Sample output:"
    head -3 "test_output.csv"
else
    echo "Verification: FAIL"
    exit 1
fi

echo "All tests passed!"
rm -f test_output.csv
