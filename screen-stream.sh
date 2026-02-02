#!/bin/bash
# Screen Streaming Script for macOS
# 使用QuickTime录屏后推流的方法

STREAM_KEY="relay_screen_test"
SRS_SERVER="rtmp://localhost/live"

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║            🐯 屏幕推流测试 - 完整指南                         ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo ""

# 检查推流进程
RUNNING_PID=$(pgrep -f "ffmpeg.*$STREAM_KEY" 2>/dev/null)
if [ -n "$RUNNING_PID" ]; then
    echo "⚠️  检测到已有推流进程 (PID: $RUNNING_PID)"
    echo "   是否停止? (y/n)"
    read -r answer
    if [ "$answer" = "y" ]; then
        kill $RUNNING_PID 2>/dev/null
        sleep 1
        echo "✅ 已停止"
    else
        echo "使用现有推流..."
        echo "🌐 播放地址: http://localhost:5173/live/a1111111-1111-1111-1111-111111111111"
        exit 0
    fi
fi

echo ""
echo "方法1: 使用测试视频推流 (推荐)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "命令:"
echo '  ffmpeg -re -stream_loop -1 -i /tmp/huya-streams/test_stream.mp4 \'
echo '      -c copy -f flv rtmp://localhost/live/relay_screen_test'
echo ""

echo "方法2: 屏幕录屏推流 (需要OBS)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "推荐使用 OBS Studio:"
echo "1. 下载: https://obsproject.com/"
echo "2. 新建场景"
echo "3. 添加'显示器捕获'源"
echo "4. 设置推流: rtmp://localhost/live"
echo "5. Stream Key: relay_screen_test"
echo "6. 开始推流"
echo ""

echo "方法3: 尝试命令行屏幕捕获"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo ""
echo "选项:"
echo "  1. 使用测试视频推流"
echo "  2. 显示OBS配置说明"
echo "  3. 退出"
echo ""
read -r choice

case $choice in
    1)
        echo ""
        echo "🚀 开始推流测试视频..."
        
        # 确保测试视频存在
        if [ ! -f "/tmp/huya-streams/test_stream.mp4" ]; then
            echo "📹 创建测试视频 (60秒)..."
            ffmpeg -f lavfi -i "testsrc=s=1280x720:r=24,format=yuv420p" \
                   -t 60 -c:v libx264 -preset ultrafast -crf 23 \
                   /tmp/huya-streams/test_stream.mp4 -y 2>/dev/null
            echo "✅ 测试视频创建完成"
        fi
        
        nohup ffmpeg -re -stream_loop -1 -i /tmp/huya-streams/test_stream.mp4 \
            -c copy -f flv "$SRS_SERVER/$STREAM_KEY" -nostdin \
            > /tmp/ffmpeg-screen.log 2>&1 &
        
        PUSH_PID=$!
        sleep 3
        
        if ps -p $PUSH_PID > /dev/null 2>&1; then
            echo "✅ 推流已启动! (PID: $PUSH_PID)"
            echo ""
            echo "🌐 请在浏览器打开:"
            echo "   http://localhost:5173/live/a1111111-1111-1111-1111-111111111111"
        else
            echo "❌ 推流启动失败"
            echo "日志:"
            tail -20 /tmp/ffmpeg-screen.log
        fi
        ;;
        
    2)
        echo ""
        echo "📺 OBS Studio 配置"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo ""
        echo "1. 下载并安装 OBS Studio"
        echo "   https://obsproject.com/download"
        echo ""
        echo "2. 打开 OBS，进入 '设置' → '推流'"
        echo ""
        echo "3. 配置:"
        echo "   服务: 自定义..."
        echo "   服务器: rtmp://localhost/live"
        echo "   串流密钥: relay_screen_test"
        echo ""
        echo "4. 添加来源:"
        echo "   + → 显示器捕获 (捕捉整个屏幕)"
        echo "   或"
        echo "   + → 窗口捕获 (捕捉特定窗口)"
        echo ""
        echo "5. 点击 '开始推流'"
        echo ""
        echo "6. 访问直播间:"
        echo "   http://localhost:5173/live/a1111111-1111-1111-1111-111111111111"
        ;;
        
    *)
        echo "退出"
        ;;
esac
