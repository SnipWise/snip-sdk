#!/bin/bash
SERVICE_URL=${SERVICE_URL:-"http://0.0.0.0:9100/api/cancel-stream-completion"}

echo "Sending cancel request to: ${SERVICE_URL}"

response=$(curl -s -w "\n%{http_code}" -X POST "${SERVICE_URL}" \
  -H "Content-Type: application/json")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ]; then
  echo "✓ Stream completion canceled successfully"
  if [ -n "$body" ]; then
    echo "Response: $body"
  fi
else
  echo "✗ Error canceling stream (HTTP $http_code)"
  if [ -n "$body" ]; then
    echo "Response: $body"
  fi
fi
