#!/bin/bash

# Test script for Mini-Dropbox file upload
# This script tests the upload functionality of storage nodes

echo "🚀 Testing Mini-Dropbox File Upload"
echo "=================================="

# Create test file
echo "Creating test file..."
echo "Hello, this is a test file for Mini-Dropbox!" > test_file.txt
echo "Test file created: test_file.txt"

# Test upload to Node 8001
echo ""
echo "📤 Testing upload to Node 8001..."
curl -X POST -F "file=@test_file.txt" http://localhost:8001/upload

echo ""
echo ""

# Test upload to Node 8002
echo "📤 Testing upload to Node 8002..."
curl -X POST -F "file=@test_file.txt" http://localhost:8002/upload

echo ""
echo ""

# Clean up test file
echo "🧹 Cleaning up test file..."
rm test_file.txt

echo ""
echo "✅ Upload testing completed!"
echo "Check the responses above for success messages and file hashes."
