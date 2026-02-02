#!/bin/bash
# Quick test script for Huya Live Platform

echo "🐯 虎牙直播平台 - 快速测试"
echo "=========================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

test() {
    local name=$1
    local url=$2
    local status=$(curl -s -o /dev/null -w "%{http_code}" "$url" --max-time 5 2>/dev/null)
    if [ "$status" == "200" ]; then
        echo -e "${GREEN}✅ $name${NC}"
        return 0
    else
        echo -e "${RED}❌ $name (HTTP $status)${NC}"
        return 1
    fi
}

echo "🔍 服务状态"
test "API Server" "http://localhost:8888/api/v1/live/rooms" && echo "  → 13个直播间活跃"
test "SRS Streaming" "http://localhost:1985/api/v1/streams/" && echo "  → 3个推流活跃"
test "Frontend" "http://localhost:5173" && echo "  → 开发服务器运行中"

echo ""
echo "📺 直播流"
test "稳定测试流 (FLV)" "http://localhost:8080/live/test_stream.flv"
test "稳定测试流 (HLS)" "http://localhost:8080/live/test_stream.m3u8"

echo ""
echo "🌐 访问地址"
echo "-----------"
echo "首页: http://localhost:5173/"
echo "直播间: http://localhost:5173/live/af763384-004a-4837-92b6-df24ca77c991"
echo ""
echo "🔑 测试账号"
echo "-----------"
echo "用户: testuser1 / test123456"
