#!/bin/bash

# Test script for Mini-Dropbox Master Node
# This script tests the master node API endpoints

echo "🎯 Testing Mini-Dropbox Master Node"
echo "=================================="

MASTER_URL="http://localhost:9000"

echo "🔍 Testing Master Node endpoints..."
echo ""

# Test 1: List all files (should be empty initially)
echo "📋 Test 1: Listing all files (should be empty initially)..."
curl -X GET "$MASTER_URL/list"
echo ""
echo ""

# Test 2: Get file by hash (should return 404 for non-existent hash)
echo "🔍 Test 2: Getting non-existent file..."
curl -X GET "$MASTER_URL/get?hash=nonexistent123"
echo ""
echo ""

# Test 3: Register a file (simulate file registration)
echo "📝 Test 3: Registering a sample file..."
curl -X POST "$MASTER_URL/register" \
  -H "Content-Type: application/json" \
  -d '{
    "hash": "test123hash456",
    "name": "sample.txt",
    "size": 1024,
    "replicas": ["8001", "8002"]
  }'
echo ""
echo ""

# Test 4: List all files again (should show the registered file)
echo "📋 Test 4: Listing all files (should show registered file)..."
curl -X GET "$MASTER_URL/list"
echo ""
echo ""

# Test 5: Get the registered file by hash
echo "🔍 Test 5: Getting registered file by hash..."
curl -X GET "$MASTER_URL/get?hash=test123hash456"
echo ""
echo ""

echo "✅ Master Node testing completed!"
echo ""
echo "📊 Summary:"
echo "- List endpoint: Shows all registered files"
echo "- Get endpoint: Retrieves file metadata by hash"
echo "- Register endpoint: Registers new file metadata"
