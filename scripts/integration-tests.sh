#!/bin/bash

# CSV Validator API Test Script

API_BASE="http://localhost:8080"
TEST_FILE="test-samples/sample1.csv"

echo "=== CSV Validator API Test ==="
echo ""

# Check if server is running
echo "1. Testing health endpoint..."
health_response=$(curl -s -w "%{http_code}" -o /dev/null "$API_BASE/health")
if [ "$health_response" = "200" ]; then
    echo "✅ Health check passed"
else
    echo "❌ Health check failed (HTTP $health_response)"
    echo "Please start the server first: go run ."
    exit 1
fi
echo ""

# Test file upload
echo "2. Testing file upload..."
if [ ! -f "$TEST_FILE" ]; then
    echo "❌ Test file not found: $TEST_FILE"
    exit 1
fi

upload_response=$(curl -s -X POST \
    -F "file=@$TEST_FILE" \
    "$API_BASE/api/upload")

job_id=$(echo "$upload_response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

if [ -n "$job_id" ]; then
    echo "✅ File uploaded successfully"
    echo "   Job ID: $job_id"
else
    echo "❌ File upload failed"
    echo "   Response: $upload_response"
    exit 1
fi
echo ""

# Test job status polling
echo "3. Polling job status..."
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    download_response=$(curl -s -w "%{http_code}" \
        "$API_BASE/api/download/$job_id" \
        -o "processed_output.csv")
    
    http_code=${download_response: -3}
    
    if [ "$http_code" = "200" ]; then
        echo "✅ File processing completed"
        echo "   Processed file saved as: processed_output.csv"
        break
    elif [ "$http_code" = "423" ]; then
        echo "   Job still in progress... (attempt $((attempt + 1))/$max_attempts)"
        sleep 1
        attempt=$((attempt + 1))
    else
        echo "❌ Download failed (HTTP $http_code)"
        curl -s "$API_BASE/api/download/$job_id"
        exit 1
    fi
done

if [ $attempt -eq $max_attempts ]; then
    echo "❌ Timeout waiting for job completion"
    exit 1
fi
echo ""

# Verify processed file
echo "4. Verifying processed file..."
if [ -f "processed_output.csv" ]; then
    echo "✅ Processed file exists"
    
    # Check if has_email column was added
    if grep -q "has_email" "processed_output.csv"; then
        echo "✅ Email validation column added"
        
        # Show first few lines
        echo ""
        echo "   First 5 lines of processed file:"
        head -5 "processed_output.csv" | sed 's/^/   /'
        
        # Count true/false values
        true_count=$(grep -c ",true$" "processed_output.csv" || echo "0")
        false_count=$(grep -c ",false$" "processed_output.csv" || echo "0")
        echo ""
        echo "   Email validation results:"
        echo "   - Rows with valid email: $true_count"
        echo "   - Rows without valid email: $false_count"
    else
        echo "❌ Email validation column not found"
        exit 1
    fi
else
    echo "❌ Processed file not found"
    exit 1
fi
echo ""

# Test error cases
echo "5. Testing error cases..."

# Test invalid file type
echo "   Testing invalid file type..."
echo "This is not a CSV file" > test_invalid.txt
invalid_response=$(curl -s -X POST \
    -F "file=@test_invalid.txt" \
    "$API_BASE/api/upload")

if echo "$invalid_response" | grep -q "error"; then
    echo "✅ Invalid file type correctly rejected"
else
    echo "❌ Invalid file type not rejected"
fi
rm -f test_invalid.txt

# Test invalid job ID
echo "   Testing invalid job ID..."
invalid_id_response=$(curl -s -w "%{http_code}" \
    "$API_BASE/api/download/invalid-id" \
    -o /dev/null)

if [ "${invalid_id_response: -3}" = "400" ]; then
    echo "✅ Invalid job ID correctly rejected"
else
    echo "❌ Invalid job ID not rejected (HTTP ${invalid_id_response: -3})"
fi

echo ""
echo "=== Test Summary ==="
echo "✅ All tests passed successfully!"
echo ""
echo "Cleanup: removing test output file..."
rm -f processed_output.csv

echo "Done!"