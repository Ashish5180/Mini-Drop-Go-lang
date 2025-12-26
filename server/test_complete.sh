#!/bin/bash

# Complete test script for Mini-Dropbox
# This script runs all tests in sequence

set -e  # Exit on error

START_TIME=$(date +%s)

echo "🚀 Mini-Dropbox Complete Testing Suite"
echo "======================================"
echo ""

# Make all test scripts executable
chmod +x test_*.sh

echo "🔧 Setting up test environment..."
echo ""

# Test 1: Health Checks
echo "🏥 Running Health Checks..."
start=$(get_ms)
./test_health.sh
end=$(get_ms)
echo "⏱️  Health check duration: $((end-start))ms"
echo ""
echo "----------------------------------------"
echo ""

# Test 2: Master Node API
echo "🎯 Testing Master Node API..."
start=$(get_ms)
./test_master.sh
end=$(get_ms)
echo "⏱️  Master API test duration: $((end-start))ms"
echo ""
echo "----------------------------------------"
echo ""

# Test 3: File Upload
echo "📤 Testing File Upload..."
start=$(get_ms)
./test_upload.sh
end=$(get_ms)
echo "⏱️  Upload test duration: $((end-start))ms"
echo ""
echo "----------------------------------------"
echo ""

# Test 4: File Retrieval (using hash from upload)
echo "🔍 Testing File Retrieval..."
echo "Note: You'll need to manually provide the hash from upload response"
echo "Example: ./test_retrieve.sh <hash_from_upload>"
echo ""

END_TIME=$(get_ms)
TOTAL_TIME=$((END_TIME-START_TIME))

echo "✅ Complete testing suite finished!"
echo ""
echo "⏱️  Total execution time: ${TOTAL_TIME}ms"
echo ""
echo "📋 Test Summary:"
echo "- ✅ Health checks: All nodes responding"
echo "- ✅ Master API: File registration and retrieval"
echo "- ✅ File upload: Storage nodes accepting files"
echo "- ℹ️  File retrieval: Run manually with hash"
echo ""
echo "💡 Next steps:"
echo "1. Check the responses above for any errors"
echo "2. Use the hash from upload response to test retrieval"
echo "3. Verify files are stored in data/ directories"
