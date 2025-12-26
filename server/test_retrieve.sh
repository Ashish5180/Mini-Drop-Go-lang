#!/bin/bash

# Test script for Mini-Dropbox file retrieval
# This script tests the file retrieval functionality

echo "🔍 Testing Mini-Dropbox File Retrieval"
echo "====================================="

# Check if hash parameter is provided
if [ -z "$1" ]; then
    echo "❌ Error: Please provide a file hash as parameter"
    echo "Usage: ./test_retrieve.sh <file_hash>"
    echo ""
    echo "Example: ./test_retrieve.sh abc123def456"
    exit 1
fi

HASH=$1

echo "🔑 Testing retrieval with hash: $HASH"
echo ""

# Test retrieval from Node 8001
echo "📥 Testing retrieval from Node 8001..."
curl -X GET "http://localhost:8001/retrieve?hash=$HASH" -o retrieved_file_8001.txt

if [ $? -eq 0 ]; then
    echo "✅ File retrieved from Node 8001: retrieved_file_8001.txt"
    echo "File content:"
    cat retrieved_file_8001.txt
else
    echo "❌ Failed to retrieve file from Node 8001"
fi

echo ""
echo ""

# Test retrieval from Node 8002
echo "📥 Testing retrieval from Node 8002..."
curl -X GET "http://localhost:8001/retrieve?hash=$HASH" -o retrieved_file_8002.txt

if [ $? -eq 0 ]; then
    echo "✅ File retrieved from Node 8002: retrieved_file_8002.txt"
    echo "File content:"
    cat retrieved_file_8002.txt
else
    echo "❌ Failed to retrieve file from Node 8002"
fi

echo ""
echo "🧹 Cleaning up retrieved files..."
rm -f retrieved_file_8001.txt retrieved_file_8002.txt

echo ""
echo "✅ Retrieval testing completed!"
