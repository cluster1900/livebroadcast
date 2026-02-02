#!/bin/bash

# Relay stream script - pulls from source and pushes to SRS

SOURCE_URL="$1"
CHANNEL_NAME="$2"
LOG_FILE="/tmp/relay_${CHANNEL_NAME}.log"

if [ -z "$SOURCE_URL" ] || [ -z "$CHANNEL_NAME" ]; then
    echo "Usage: $0 <source_url> <channel_name>"
    exit 1
fi

# Kill existing relay for this channel
pkill -f "ffmpeg.*${CHANNEL_NAME}" 2>/dev/null
sleep 1

# Start relay
echo "Starting relay: $SOURCE_URL -> rtmp://localhost/live/$CHANNEL_NAME" >> "$LOG_FILE"
ffmpeg -re -i "$SOURCE_URL" -c copy -f flv "rtmp://localhost/live/$CHANNEL_NAME" -nostdin >> "$LOG_FILE" 2>&1 &

echo "Relay started with PID: $!"
