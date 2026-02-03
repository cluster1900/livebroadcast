#!/bin/bash
# Huya Live - Stream Push Script
# 推送流到直播间对应的channel_name

STREAM_SERVER="rtmp://localhost/live"

echo "🚀 Starting push streams..."

# Clean up old streams
pkill -f "ffmpeg.*live_" 2>/dev/null
pkill -f "ffmpeg.*relay_" 2>/dev/null
sleep 1

# Get channel names from database that need streams
# For demo, we'll push to the test streams that exist in DB

# Push to relay_c168a5b77c81979f (稳定测试流)
echo "Starting stable test stream..."
nohup ffmpeg -re -f lavfi -i testsrc=s=1280x720:r=30 \
    -c:v libx264 -preset ultrafast -tune zerolatency -c:a aac -b:a 128k \
    -f flv "${STREAM_SERVER}/relay_c168a5b77c81979f" -nostdin > /tmp/ffmpeg-stable.log 2>&1 &

# Push to relay_433e022aa907e572 (本地视频测试v2)
echo "Starting local video test v2 stream..."
nohup ffmpeg -re -f lavfi -i testsrc=s=854x480:r=25 \
    -c:v libx264 -preset ultrafast -tune zerolatency -c:a aac -b:a 96k \
    -f flv "${STREAM_SERVER}/relay_433e022aa907e572" -nostdin > /tmp/ffmpeg-local2.log 2>&1 &

# Push to relay_screen_test (屏幕推流测试)
echo "Starting screen relay stream..."
nohup ffmpeg -f avfoundation -capture_cursor 1 -i 3 \
    -c:v libx264 -preset ultrafast -tune zerolatency -pix_fmt yuv420p -b:v 2000k \
    -f flv "${STREAM_SERVER}/relay_screen_test" -nostdin > /tmp/ffmpeg-screen.log 2>&1 &

# Push to test_stream (用于普通测试直播间)
echo "Starting test stream..."
nohup ffmpeg -re -f lavfi -i testsrc=s=640x360:r=24 \
    -c:v libx264 -preset ultrafast -tune zerolatency -c:a aac -b:a 64k \
    -f flv "${STREAM_SERVER}/test_stream" -nostdin > /tmp/ffmpeg-test.log 2>&1 &

# Push to hd_stream (用于高清测试直播间)
echo "Starting HD stream..."
nohup ffmpeg -re -f lavfi -i testsrc=s=1280x720:r=30 \
    -c:v libx264 -preset ultrafast -tune zerolatency -c:a aac -b:a 128k \
    -f flv "${STREAM_SERVER}/hd_stream" -nostdin > /tmp/ffmpeg-hd.log 2>&1 &

echo "✅ Push streams started"
sleep 2
echo ""
echo "📺 Active streams:"
ps aux | grep ffmpeg | grep -v grep | awk '{print $NF}' | while read stream; do
    echo "  • $stream"
done
