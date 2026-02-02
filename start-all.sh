#!/bin/bash
# Huya Live Platform - Complete Startup Script
# Starts all services, HTTP server for test videos, and push streams

set -e

echo "🐯 虎牙直播平台 - 完整启动脚本"
echo "================================"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Function to check command status
check_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
    fi
}

# Step 1: Start Docker services
echo "🔧 Step 1: 启动 Docker 服务..."
docker start huya_live-postgres-1 huya_live-redis-1 huya_live-centrifugo-1 srs 2>/dev/null || true
sleep 3
docker ps --format "table {{.Names}}\t{{.Status}}" | grep -E "huya|srs" | head -5
echo ""

# Step 2: Start HTTP server for test videos
echo "🌐 Step 2: 启动测试视频 HTTP 服务器..."
mkdir -p /tmp/huya-streams

# Create test video if not exists
if [ ! -f "/tmp/huya-streams/test_stream.mp4" ]; then
    echo "   创建测试视频..."
    ffmpeg -f lavfi -i "testsrc=s=640x360:r=24,format=yuv420p" \
           -t 300 -c:v libx264 -preset ultrafast -crf 23 \
           /tmp/huya-streams/test_stream.mp4 -y 2>/dev/null &
    echo "   视频生成中..."
fi

# Start HTTP server
pkill -f "python3.*http.server 7777" 2>/dev/null || true
sleep 1
nohup python3 -m http.server 7777 -d /tmp/huya-streams > /tmp/http-7777.log 2>&1 &
sleep 2

# Verify HTTP server
if curl -m 3 http://localhost:7777/test_stream.mp4 > /dev/null 2>&1; then
    check_status 0 "HTTP 服务器 (端口 7777)"
else
    check_status 1 "HTTP 服务器 (端口 7777)"
fi
echo ""

# Step 3: Start API server
echo "🚀 Step 3: 启动 API 服务器..."
pkill -f "huya.*api\|./server" 2>/dev/null || true
sleep 1
cd /Users/hawkwu/Desktop/huya_live/api
nohup env DB_HOST=localhost DB_PASSWORD=huya_live_secret REDIS_ADDR=localhost:6379 \
    ./server > /tmp/huya-api.log 2>&1 &
sleep 3

# Verify API
if curl -m 5 http://localhost:8888/api/v1/live/rooms > /dev/null 2>&1; then
    check_status 0 "API 服务器 (端口 8888)"
else
    check_status 1 "API 服务器 (端口 8888)"
fi
echo ""

# Step 4: Start push streams
echo "📺 Step 4: 启动推流..."

# Clean up old streams
pkill -f "ffmpeg.*-f flv rtmp://localhost/live" 2>/dev/null || true
sleep 1

STREAM_SERVER="rtmp://localhost/live"

# Start test streams
echo "   启动基础测试流..."

nohup ffmpeg -re -stream_loop -1 -i /tmp/huya-streams/test_stream.mp4 \
    -c copy -f flv "${STREAM_SERVER}/local_test_stream" -nostdin \
    > /tmp/ffmpeg-local.log 2>&1 &

nohup ffmpeg -re -f lavfi -i "testsrc=s=640x360:r=24" \
    -c:v libx264 -preset ultrafast -tune zerolatency -c:a aac -b:a 64k \
    -f flv "${STREAM_SERVER}/test_stream" -nostdin \
    > /tmp/ffmpeg-test.log 2>&1 &

nohup ffmpeg -re -f lavfi -i "testsrc=s=1280x720:r=30" \
    -c:v libx264 -preset ultrafast -tune zerolatency -c:a aac -b:a 128k \
    -f flv "${STREAM_SERVER}/hd_stream" -nostdin \
    > /tmp/ffmpeg-hd.log 2>&1 &

# Start relay streams from database
echo "   启动数据库中的中转流..."

RELAY_STREAMS=$(docker exec huya_live-postgres-1 psql -U huya_live -d huya_live -t -A -c "
    SELECT channel_name, source_url FROM relay_streams 
    WHERE status='running' AND source_url LIKE 'http://localhost%';
" 2>/dev/null)

if [ -n "$RELAY_STREAMS" ]; then
    echo "$RELAY_STREAMS" | while IFS='|' read -r channel_name source_url; do
        if [ -z "$channel_name" ] || [ -z "$source_url" ]; then
            continue
        fi
        
        # Use local file directly instead of HTTP
        local_file=$(echo "$source_url" | sed 's|http://localhost:7777/|/tmp/huya-streams/|g')
        
        if [ -f "$local_file" ]; then
            nohup ffmpeg -re -stream_loop -1 -i "$local_file" \
                -c copy -f flv "${STREAM_SERVER}/${channel_name}" -nostdin \
                > "/tmp/ffmpeg-${channel_name}.log" 2>&1 &
            echo "   📺 $channel_name (local file)"
        else
            nohup ffmpeg -re -stream_loop -1 -i "$source_url" \
                -c copy -f flv "${STREAM_SERVER}/${channel_name}" -nostdin \
                > "/tmp/ffmpeg-${channel_name}.log" 2>&1 &
            echo "   📺 $channel_name (HTTP)"
        fi
    done
fi

sleep 3
echo ""

# Step 5: Verify streams
echo "🔍 Step 5: 验证推流状态..."
sleep 2

# Quick stream test
STREAM_TEST=$(curl -m 5 http://localhost:8080/live/test_stream.flv -o /dev/null -w '%{http_code}' 2>/dev/null || echo "000")
if [ "$STREAM_TEST" == "200" ]; then
    check_status 0 "测试流 (test_stream)"
else
    check_status 1 "测试流 (test_stream)"
fi

STREAM_TEST2=$(curl -m 5 http://localhost:8080/live/local_test_stream.flv -o /dev/null -w '%{http_code}' 2>/dev/null || echo "000")
if [ "$STREAM_TEST2" == "200" ]; then
    check_status 0 "本地流 (local_test_stream)"
else
    check_status 1 "本地流 (local_test_stream)"
fi

echo ""
echo "📊 当前推流进程:"
ps aux | grep ffmpeg | grep -E "test_stream|local_test|hd_stream|relay_" | grep -v grep | wc -l
echo ""

# Summary
echo "================================"
echo -e "${GREEN}🎉 启动完成！${NC}"
echo "================================"
echo ""
echo "🌐 访问地址:"
echo "   前端: http://localhost:5173"
echo "   直播间: http://localhost:5173/live/af763384-004a-4837-92b6-df24ca77c991"
echo ""
echo "🔑 测试账号:"
echo "   用户: testuser1 / test123456"
echo "   主播: testuser2 / test123456"
echo ""
echo "📝 日志位置:"
echo "   API: /tmp/huya-api.log"
echo "   HTTP: /tmp/http-7777.log"
echo "   FFmpeg: /tmp/ffmpeg-*.log"
