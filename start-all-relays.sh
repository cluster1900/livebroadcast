#!/bin/bash
# Huya Live - Start All Relay Streams
# Reads relay streams from database and starts FFmpeg processes

RELAY_SERVER="http://localhost:8080"

echo "🚀 启动所有中转流..."
echo "===================="

# Kill existing relay streams
pkill -f "ffmpeg.*relay_" 2>/dev/null
sleep 1

# Get relay streams from database
RELAY_STREAMS=$(docker exec huya_live-postgres-1 psql -U huya_live -d huya_live -t -A -c "SELECT channel_name, source_url FROM relay_streams WHERE status='running';" 2>/dev/null)

if [ -z "$RELAY_STREAMS" ]; then
    echo "⚠️  未找到运行中的中转流配置"
    exit 1
fi

echo "$RELAY_STREAMS" | while IFS='|' read -r channel_name source_url; do
    if [ -z "$channel_name" ]; then
        continue
    fi
    
    echo "📺 启动中转流: $channel_name"
    echo "   源: $source_url"
    
    # Determine input options based on source URL
    if [[ "$source_url" == *.m3u8* ]] || [[ "$source_url" == http://* ]] || [[ "$source_url" == https://* ]]; then
        # HTTP/HLS source
        nohup ffmpeg -re -i "$source_url" \
            -c copy -f flv "rtmp://localhost/live/$channel_name" -nostdin \
            > "/tmp/ffmpeg-$channel_name.log" 2>&1 &
    elif [[ "$source_url" == rtmp://* ]]; then
        # RTMP source
        nohup ffmpeg -re -i "$source_url" \
            -c copy -f flv "rtmp://localhost/live/$channel_name" -nostdin \
            > "/tmp/ffmpeg-$channel_name.log" 2>&1 &
    else
        echo "   ⚠️  未知源类型: $source_url"
    fi
    
    sleep 0.5
done

echo ""
echo "⏳ 等待流稳定..."
sleep 3

echo ""
echo "📊 当前推流进程:"
ps aux | grep ffmpeg | grep -v grep | awk '{print $NF}' | head -10

echo ""
echo "✅ 中转流启动完成！"
