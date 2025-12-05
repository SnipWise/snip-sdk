#!/bin/bash
SERVICE_URL=${SERVICE_URL:-"http://0.0.0.0:9200/api/chat-stream"}

read -r -d '' USER_CONTENT <<- EOM
Hello, my name is Bob Morane. Tell me a story about your adventures.
EOM

read -r -d '' DATA <<- EOM
{
  "data": {
    "message":"${USER_CONTENT}"
  }
}
EOM

# Remove newlines from DATA 
DATA=$(echo ${DATA} | tr -d '\n')

echo "Using DATA: ${DATA}"
echo -e "\n"


callback() {
  echo -ne "$1" 
}

unescape_quotes() {
    local str="$1"
    str="${str//\\\"/\"}"  # Replace \" by "
    echo "$str"
}

remove_quotes() {
    local str="$1"
    str="${str%\"}"   # remove " at the end
    str="${str#\"}"   # remove " at start
    echo "$str"
}

curl --no-buffer --silent ${SERVICE_URL} \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
    -d "${DATA}" \
  | while IFS= read -r line; do
    if [[ $line == data:* ]]; then
      json_data="${line#data: }"

      # Extract from nested message.response structure (Genkit format)
      content_chunk=$(echo "$json_data" | jq '.message.response // empty' 2>/dev/null)

      # Check for finish_reason to detect end of stream
      finish_reason=$(echo "$json_data" | jq -r '.message.finish_reason // empty' 2>/dev/null)

      if [[ -n "$content_chunk" && "$content_chunk" != '""' ]]; then
        result=$(remove_quotes "$content_chunk")
        clean_result=$(unescape_quotes "$result")
        callback "$clean_result"
      fi

      # Display finish reason if present
      if [[ -n "$finish_reason" && "$finish_reason" != "null" ]]; then
        echo ""
        echo "[Stream completed - Finish reason: $finish_reason]"
      fi
    fi
  done

echo ""