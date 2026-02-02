#!/bin/bash
# Huya Live - Stream Push Script

STREAM_SERVER="rtmp://localhost/live"

# Check if test video exists
if [ ! -f "/tmp/test_stream.mp4" ]; then
    echo "Downloading test video..."
    curl -L -o /tmp/test_stream.mp4 "https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8" 2>/dev/null || \
    curl -L -o /tmp/test_stream.mp4 "https://media.w3.org/2010/05/sintel/trailer.mp4" 2>/dev/null || \
    ffmpeg -f lavfi -i testsrc=s=640x360:r=24 -c:v libx264 -preset ultrafast -t 60 -f mp4 /tmp/test_stream.mp4 2>/dev/null
fi

echo "🚀 Starting push streams..."

# Clean up old streams
pkill -f "ffmpeg.*live_" 2>/dev/null
sleep 1

# Push streams in background
nohup ffmpeg -re -stream_loop -1 -i /tmp/test_stream.mp4 \
    -c copy -f flv "${STREAM_SERVER}/local_test_stream" -nostdin > /tmp/ffmpeg-local.log 2>&1 &

nohup ffmpeg -re -f lavfi -i testsrc=s=640x360:r=24 \
    -c:v libx264 -preset ultrafast -tune zerolatency -c:a aac -b:a 64k \
    -f flv "${STREAM_SERVER}/test_stream" -nostdin > /tmp/ffmpeg-test.log 2>&1 &

nohup ffmpeg -re -f lavfi -i testsrc=s=1280x720:r=30 \
    -c:v libx264 -preset ultrafast -tune zerolatency -c:a aac -b:a 128k \
    -f flv "${STREAM_SERVER}/hd_stream" -nostdin > /tmp/ffmpeg-hd.log 2>&1 &

echo "✅ Push streams started"
sleep 2
ps aux | grep ffmpeg | grep -v grep
