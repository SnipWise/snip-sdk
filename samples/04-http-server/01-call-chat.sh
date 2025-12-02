#!/bin/bash
SERVICE_URL=${SERVICE_URL:-"http://0.0.0.0:9100/api/chat"}

read -r -d '' USER_CONTENT <<- EOM
Hello, who are you? Please introduce yourself in a concise manner.
EOM

read -r -d '' DATA <<- EOM
{
  "data": {
    "message":"${USER_CONTENT}"
  }
}
EOM

curl -X POST ${SERVICE_URL} \
  -H "Content-Type: application/json" \
  -d "${DATA}"