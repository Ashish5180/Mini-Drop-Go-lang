#!/bin/bash
set -euo pipefail  # Better error handling

echo "🏥 Testing Mini-Dropbox Health Checks"
echo "===================================="

echo "🔍 Testing health endpoints..."
echo ""

# Test Master Node health (port 9000)
echo "🎯 Testing Master Node health (port 9000)..."
if curl -f -s -m 2 -X GET "http://localhost:9000/health" 2>/dev/null; then
    echo " - ✅ Master Node responding"
else
    echo " - ⚠️  Master Node health endpoint not implemented"
fi
echo ""
echo ""

# Test Storage Node 8001 health
echo "📦 Testing Storage Node 8001 health..."
if curl -f -s -m 2 -X GET "http://localhost:8001/health"; then
    echo ""
    echo " - ✅ Storage Node 8001 healthy"
else
    echo " - ❌ Storage Node 8001 not responding"
fi
echo ""
echo ""

# Test Storage Node 8002 health
echo "📦 Testing Storage Node 8002 health..."
if curl -f -s -m 2 -X GET "http://localhost:8002/health"; then
    echo ""
    echo " - ✅ Storage Node 8002 healthy"
else
    echo " - ❌ Storage Node 8002 not responding"
fi
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
