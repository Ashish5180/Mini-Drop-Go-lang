#!/bin/bash

# Test script for Mini-Dropbox Health Checks
# This script tests the health endpoints of all nodes

echo "🏥 Testing Mini-Dropbox Health Checks"
echo "===================================="

echo "🔍 Testing health endpoints..."
echo ""

# Test Master Node health (port 9000)
echo "🎯 Testing Master Node health (port 9000)..."
curl -X GET "http://localhost:9000/health" 2>/dev/null || echo "❌ Master Node not responding"
echo ""
echo ""

# Test Storage Node 8001 health
echo "📦 Testing Storage Node 8001 health..."
curl -X GET "http://localhost:8001/health"
echo ""
echo ""

# Test Storage Node 8002 health
echo "📦 Testing Storage Node 8002 health..."
curl -X GET "http://localhost:8002/health"
echo ""
echo ""

echo "✅ Health check testing completed!"
echo ""
echo "📊 Expected Results:"
echo "- Master Node: May not have health endpoint (not implemented)"
echo "- Storage Nodes: Should return 'OK' status"
echo ""
echo "💡 If any node is not responding, make sure it's running:"
echo "   go run cmd/main.go"
