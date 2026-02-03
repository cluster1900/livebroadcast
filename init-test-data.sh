#!/bin/bash

set -e

API_URL="${API_URL:-http://localhost:8888}"

echo "🐯 初始化测试数据..."

# 1. 创建测试用户
echo "创建测试用户..."

# 用户1 - 普通用户
curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser1", "password": "test123456", "nickname": "测试用户1", "email": "test1@example.com"}' | jq .

# 用户2 - 主播
curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser2", "password": "test123456", "nickname": "测试主播", "email": "test2@example.com"}' | jq .

# 用户3 - 另一个主播
curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser3", "password": "test123456", "nickname": "主播小姐姐", "email": "test3@example.com"}' | jq .

# 用户4 - 管理员
curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123", "nickname": "管理员", "email": "admin@example.com"}' | jq .

echo ""
echo "用户创建完成！登录获取token..."

# 2. 获取token
TOKEN=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser1", "password": "test123456"}' | jq -r '.data.access_token')

USER1_ID=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser1", "password": "test123456"}' | jq -r '.data.user.id')

USER2_ID=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser2", "password": "test123456"}' | jq -r '.data.user.id')

USER3_ID=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser3", "password": "test123456"}' | jq -r '.data.user.id')

ADMIN_ID=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}' | jq -r '.data.user.id')

echo "USER1_ID: $USER1_ID"
echo "USER2_ID: $USER2_ID"
echo "USER3_ID: $USER3_ID"
echo "ADMIN_ID: $ADMIN_ID"

# 3. 为主播用户申请主播资格
echo ""
echo "申请主播资格..."

curl -s -X POST "$API_URL/api/v1/streamers/apply" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"channel_name": "test_streamer"}' | jq .

# 4. 创建直播间
echo ""
echo "创建直播间..."

ROOM1=$(curl -s -X POST "$API_URL/api/v1/rooms" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "游戏直播：英雄联盟", "category": "游戏"}' | jq -r '.data.id')

ROOM2=$(curl -s -X POST "$API_URL/api/v1/rooms" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "唱歌直播：流行歌曲", "category": "音乐"}' | jq -r '.data.id')

ROOM3=$(curl -s -X POST "$API_URL/api/v1/rooms" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "户外直播：城市探险", "category": "户外"}' | jq -r '.data.id')

echo "ROOM1: $ROOM1"
echo "ROOM2: $ROOM2"
echo "ROOM3: $ROOM3"

# 5. 创建直播预告
echo ""
echo "创建直播预告..."

TOMORROW=$(date -d "+1 day" +%Y-%m-%dT%H:00:00Z)
NEXTWEEK=$(date -d "+7 day" +%Y-%m-%dT%H:00:00Z)

curl -s -X POST "$API_URL/api/v1/schedules" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"title\": \"明天晚上8点准时开播\", \"description\": \"精彩内容不容错过\", \"category\": \"娱乐\", \"start_time\": \"$TOMORROW\"}" | jq .

curl -s -X POST "$API_URL/api/v1/schedules" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"title\": \"一周特别直播\", \"description\": \"感谢大家的支持\", \"category\": \"游戏\", \"start_time\": \"$NEXTWEEK\"}" | jq .

# 6. 关注主播
echo ""
echo "关注主播..."

curl -s -X POST "$API_URL/api/v1/social/follow" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"streamer_id\": \"$USER2_ID\"}" | jq .

# 7. 发送私信
echo ""
echo "发送私信..."

curl -s -X POST "$API_URL/api/v1/messages/send" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"receiver_id\": \"$USER2_ID\", \"content\": \"主播你好，我是你的粉丝！\"}" | jq .

curl -s -X POST "$API_URL/api/v1/messages/send" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"receiver_id\": \"$USER3_ID\", \"content\": \"欢迎来到直播平台！\"}" | jq .

# 8. 充值虎牙币
echo ""
echo "充值虎牙币..."

curl -s -X POST "$API_URL/api/v1/wallet/recharge" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"amount": 1000}' | jq .

# 9. 发送礼物
echo ""
echo "发送礼物..."

curl -s -X POST "$API_URL/api/v1/gifts/send" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"room_id\": \"$ROOM1\", \"gift_id\": 1, \"count\": 10}" | jq .

# 10. 发送弹幕
echo ""
echo "发送弹幕..."

curl -s -X POST "$API_URL/api/v1/danmu/send" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"room_id\": \"$ROOM1\", \"content\": \"主播66666\"}" | jq .

# 11. 点赞直播间
echo ""
echo "点赞直播间..."

curl -s -X POST "$API_URL/api/v1/likes/rooms/$ROOM1" \
  -H "Authorization: Bearer $TOKEN" | jq .

# 12. 添加观看历史
echo ""
echo "添加观看历史..."

curl -s -X POST "$API_URL/api/v1/history/watch" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"room_id\": \"$ROOM1\", \"watch_duration\": 3600}" | jq .

# 13. 创建举报
echo ""
echo "创建举报测试..."

curl -s -X POST "$API_URL/api/v1/reports" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"reported_id\": \"$USER2_ID\", \"type\": \"spam\", \"reason\": \"测试举报功能\"}" | jq .

# 14. 创建敏感词
echo ""
echo "创建敏感词（管理员操作）..."

# 先登录管理员
ADMIN_TOKEN=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}' | jq -r '.data.access_token')

curl -s -X POST "$API_URL/api/v1/admin/sensitive-words" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"word": "垃圾广告", "type": "blacklist", "severity": "high"}' | jq .

curl -s -X POST "$API_URL/api/v1/admin/sensitive-words" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"word": "恶意灌水", "type": "blacklist", "severity": "medium"}' | jq .

# 15. 创建礼物
echo ""
echo "创建礼物（管理员操作）..."

curl -s -X POST "$API_URL/api/v1/admin/gifts" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "豪华跑车", "coin_price": 2000, "icon_url": "/gifts/car.png", "animation_type": "lottie", "sort_order": 10}' | jq .

curl -s -X POST "$API_URL/api/v1/admin/gifts" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "浪漫玫瑰", "coin_price": 99, "icon_url": "/gifts/rose.png", "animation_type": "css", "sort_order": 9}' | jq .

# 16. 更新用户等级（模拟）
echo ""
echo "更新用户等级..."

curl -s -X PUT "$API_URL/api/v1/user/profile" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"nickname": "升级用户", "avatar_url": "https://api.dicebear.com/7.x/avataaars/svg?seed=testuser1"}' | jq .

# 17. 获取所有数据
echo ""
echo "========== 数据汇总 =========="

echo ""
echo "直播间列表:"
curl -s "$API_URL/api/v1/live/rooms" | jq '.data | length' && \
curl -s "$API_URL/api/v1/live/rooms" | jq -r '.data[] | "- \(.title) (\(.status))"' | head -10

echo ""
echo "用户关注:"
curl -s "$API_URL/api/v1/social/followings" \
  -H "Authorization: Bearer $TOKEN" | jq '.data | length'

echo ""
echo "私信会话:"
curl -s "$API_URL/api/v1/messages/conversations" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""
echo "观看历史:"
curl -s "$API_URL/api/v1/history/watch" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""
echo "礼物背包:"
curl -s "$API_URL/api/v1/inventory/gifts" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""
echo "直播预告:"
curl -s "$API_URL/api/v1/extra/schedules/upcoming" | jq '.'

echo ""
echo "通知列表:"
curl -s "$API_URL/api/v1/notifications" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""
echo "我的举报:"
curl -s "$API_URL/api/v1/reports/my" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""
echo "========== 初始化完成 =========="
echo "测试账号:"
echo "  用户1: testuser1 / test123456"
echo "  主播1: testuser2 / test123456"
echo "  主播2: testuser3 / test123456"
echo "  管理员: admin / admin123"
