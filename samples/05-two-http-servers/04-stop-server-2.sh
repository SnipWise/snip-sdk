#!/bin/bash
SERVICE_URL=${SERVICE_URL:-"http://0.0.0.0:9200/server/shutdown"}

# Test script to shutdown the server via HTTP endpoint

echo "Sending shutdown request to the server..."
curl -X POST ${SERVICE_URL} \
  -H "Content-Type: application/json"

echo ""
echo "Server should shutdown gracefully."
