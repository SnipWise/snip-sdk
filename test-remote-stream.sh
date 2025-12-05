#!/bin/bash

# Test script to check if the remote agent stream endpoint works

echo "Testing remote agent streaming..."
echo "Make sure the server is running on port 9100"
echo ""

curl -N -X POST http://0.0.0.0:9100/api/chat-stream \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{"data":{"message":"Count from 1 to 5"}}' \
  2>&1

echo ""
echo "Done"
