#!/bin/bash

# Huya Live Platform - Comprehensive Test Script
# Tests all major functionality of the platform

echo "🐯 虎牙直播平台 - 功能测试脚本"
echo "================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test function
test_endpoint() {
    local name=$1
    local url=$2
    local expected_status=$3
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null)
    
    if [ "$response" == "$expected_status" ]; then
        echo -e "${GREEN}✅ $name${NC} - HTTP $response"
        return 0
    else
        echo -e "${RED}❌ $name${NC} - Expected $expected_status, got HTTP $response"
        return 1
    fi
}

# Test results
passed=0
failed=0

echo "🔍 服务健康检查"
echo "----------------"

# Test services
test_endpoint "API Server" "http://localhost:8888/api/v1/health" "200" || ((failed++))
test_endpoint "SRS API" "http://localhost:1985/api/v1/streams/" "200" || ((failed++))
test_endpoint "Centrifugo" "http://localhost:8000/health" "200" || ((failed++))

echo ""
echo "📺 直播流测试"
echo "-------------"

# Test streams
test_endpoint "稳定测试流 (FLV)" "http://localhost:8080/live/relay_c168a5b77c81979f.flv" "200" || ((failed++))
test_endpoint "稳定测试流 (HLS)" "http://localhost:8080/live/relay_c168a5b77c81979f.m3u8" "200" || ((failed++))
test_endpoint "本地视频流 (FLV)" "http://localhost:8080/live/local_test_stream.flv" "200" || ((failed++))
test_endpoint "本地视频流 (HLS)" "http://localhost:8080/live/local_test_stream.m3u8" "200" || ((failed++))

echo ""
echo "🔗 API 端点测试"
echo "---------------"

test_endpoint "直播间列表" "http://localhost:8888/api/v1/live/rooms" "200" || ((failed++))
test_endpoint "礼物列表" "http://localhost:8888/api/v1/gifts" "200" || ((failed++))
test_endpoint "排行榜" "http://localhost:8888/api/v1/leaderboard/global" "200" || ((failed++))

echo ""
echo "📊 测试统计"
echo "-----------"

total=$((passed + failed))
if [ $failed -eq 0 ]; then
    echo -e "${GREEN}全部通过！$total/$total 项测试成功${NC}"
else
    echo -e "${YELLOW}部分测试失败：$passed/$total 通过，$failed 失败${NC}"
fi

echo ""
echo "🌐 前端访问地址"
echo "---------------"
echo "首页: http://localhost:5173/"
echo "直播间: http://localhost:5173/live/af763384-004a-4837-92b6-df24ca77c991"
echo ""
echo "🔑 测试账号"
echo "-----------"
echo "用户: testuser1 / test123456"
echo "主播: testuser2 / test123456"
