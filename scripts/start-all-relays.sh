#!/bin/bash

# 停止所有现有的推流
pkill -f "ffmpeg.*relay_" 2>/dev/null
sleep 1

echo "=== 为所有中转流创建推流 ==="

# 读取中转流列表并创建推流
declare -A relay_map=(
    ["relay_4707f38a9405a859"]="https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
    ["relay_5ae499dea409ab20"]="https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
    ["relay_cec8a445edc04a88"]="https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
    ["relay_7f6dbff5950eee00"]="https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
    ["relay_62190e68d9b72670"]="https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
    ["relay_41512c9697a86759"]="https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
    ["relay_468cd076222f78f3"]="https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
    ["relay_9a9a6fc211ba5b0f"]="https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
    ["relay_c3351ce4542d0bf2"]="https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
    ["relay_e93b5a5ee42e68e6"]="https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
    ["relay_6fb37b62a8179454"]="https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
)

for channel in "${!relay_map[@]}"; do
    source_url="${relay_map[$channel]}"
    echo "启动推流: $channel <- $source_url"
    nohup ffmpeg -re -i "$source_url" -c copy -f flv "rtmp://localhost/live/$channel" -nostdin > /tmp/relay_${channel}.log 2>&1 &
    echo "  PID: $!"
done

echo ""
echo "=== 验证推流状态 ==="
sleep 3
curl -sL "http://localhost:1985/api/v1/streams" 2>/dev/null | python3 -c "
import json,sys
d=json.load(sys.stdin)
streams=d.get('streams',[])
print(f'SRS上有 {len(streams)} 个流:')
for s in streams:
    name=s.get('name','unknown')
    clients=s.get('clients',0)
    print(f'  - {name}: {clients} clients')
"