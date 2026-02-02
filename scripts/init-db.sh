#!/bin/bash
set -e

echo "Initializing database..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- 创建扩展
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    -- 用户表
    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        username VARCHAR(50) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        nickname VARCHAR(100),
        avatar_url TEXT,
        phone VARCHAR(20) UNIQUE,
        email VARCHAR(100) UNIQUE,
        level INT DEFAULT 1,
        exp BIGINT DEFAULT 0,
        coin_balance INT DEFAULT 0,
        status VARCHAR(20) DEFAULT 'active',
        last_login_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ DEFAULT NOW(),
        updated_at TIMESTAMPTZ DEFAULT NOW()
    );

    -- 主播表
    CREATE TABLE IF NOT EXISTS streamers (
        user_id UUID PRIMARY KEY REFERENCES users(id),
        stream_key VARCHAR(128) UNIQUE NOT NULL,
        stream_key_expire_at TIMESTAMPTZ,
        rtmp_url TEXT,
        status VARCHAR(20) DEFAULT 'offline',
        is_verified BOOLEAN DEFAULT false,
        total_revenue BIGINT DEFAULT 0,
        follower_count INT DEFAULT 0,
        total_live_duration INT DEFAULT 0,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );

    -- 直播间表
    CREATE TABLE IF NOT EXISTS live_rooms (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        streamer_id UUID NOT NULL REFERENCES users(id),
        title VARCHAR(200) NOT NULL,
        category VARCHAR(50),
        cover_url TEXT,
        channel_name VARCHAR(100) UNIQUE NOT NULL,
        status VARCHAR(20) DEFAULT 'ended',
        start_at TIMESTAMPTZ,
        end_at TIMESTAMPTZ,
        peak_online INT DEFAULT 0,
        total_views INT DEFAULT 0,
        record_url TEXT,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );

    -- 礼物表
    CREATE TABLE IF NOT EXISTS gifts (
        id SERIAL PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        coin_price INT NOT NULL CHECK (coin_price > 0),
        icon_url TEXT NOT NULL,
        animation_type VARCHAR(20),
        animation_url TEXT,
        min_level_required INT DEFAULT 1,
        is_active BOOLEAN DEFAULT true,
        sort_order INT DEFAULT 0,
        category VARCHAR(20) DEFAULT 'normal',
        created_at TIMESTAMPTZ DEFAULT NOW()
    );

    -- 礼物交易表
    CREATE TABLE IF NOT EXISTS gift_transactions (
        id BIGSERIAL PRIMARY KEY,
        sender_id UUID NOT NULL REFERENCES users(id),
        receiver_id UUID NOT NULL REFERENCES users(id),
        room_id UUID NOT NULL REFERENCES live_rooms(id),
        gift_id INT NOT NULL REFERENCES gifts(id),
        gift_count INT DEFAULT 1,
        coin_amount INT NOT NULL,
        loyalty_points_gained BIGINT,
        user_level_at_send INT,
        bonus_multiplier DECIMAL(5,2) DEFAULT 1.0,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );

    -- 货币交易表
    CREATE TABLE IF NOT EXISTS coin_transactions (
        id BIGSERIAL PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users(id),
        amount INT NOT NULL,
        balance_after INT NOT NULL,
        type VARCHAR(20) NOT NULL,
        related_id BIGINT,
        description TEXT,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );

    -- 粉丝关系表
    CREATE TABLE IF NOT EXISTS fan_relations (
        user_id UUID NOT NULL REFERENCES users(id),
        streamer_id UUID NOT NULL REFERENCES users(id),
        fan_level INT DEFAULT 1,
        loyalty_points BIGINT DEFAULT 0,
        badge_name VARCHAR(20),
        badge_worn BOOLEAN DEFAULT true,
        total_gift_amount BIGINT DEFAULT 0,
        followed_at TIMESTAMPTZ DEFAULT NOW(),
        last_gift_at TIMESTAMPTZ,
        PRIMARY KEY (user_id, streamer_id)
    );

    -- 等级配置表
    CREATE TABLE IF NOT EXISTS level_config (
        level INT PRIMARY KEY,
        exp_required BIGINT NOT NULL,
        loyalty_points_required BIGINT,
        bonus_multiplier DECIMAL(5,2) DEFAULT 1.0,
        level_name VARCHAR(50),
        icon_url TEXT,
        color VARCHAR(7)
    );

    -- 敏感词表
    CREATE TABLE IF NOT EXISTS sensitive_words (
        id SERIAL PRIMARY KEY,
        word VARCHAR(50) NOT NULL,
        type VARCHAR(20) DEFAULT 'blacklist',
        severity VARCHAR(10) DEFAULT 'medium',
        is_active BOOLEAN DEFAULT true,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );

    -- 系统配置表
    CREATE TABLE IF NOT EXISTS system_config (
        key VARCHAR(50) PRIMARY KEY,
        value TEXT NOT NULL,
        description TEXT,
        updated_at TIMESTAMPTZ DEFAULT NOW()
    );

    -- 创建索引
    CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
    CREATE INDEX IF NOT EXISTS idx_users_level ON users(level DESC);
    CREATE INDEX IF NOT EXISTS idx_streamers_status ON streamers(status);
    CREATE INDEX IF NOT EXISTS idx_live_rooms_status ON live_rooms(status, start_at DESC);
    CREATE INDEX IF NOT EXISTS idx_live_rooms_streamer ON live_rooms(streamer_id, status);
    CREATE INDEX IF NOT EXISTS idx_gifts_active ON gifts(is_active, sort_order);
    CREATE INDEX IF NOT EXISTS idx_gift_transactions_room ON gift_transactions(room_id, created_at DESC);
    CREATE INDEX IF NOT EXISTS idx_coin_transactions_user ON coin_transactions(user_id, created_at DESC);
    CREATE INDEX IF NOT EXISTS idx_fan_relations_streamer ON fan_relations(streamer_id, loyalty_points DESC);
    CREATE INDEX IF NOT EXISTS idx_sensitive_words_active ON sensitive_words(is_active);

    -- 插入初始数据
    INSERT INTO level_config (level, exp_required, loyalty_points_required, bonus_multiplier, level_name, color) VALUES
        (1, 0, 0, 1.0, '萌新', '#999999'),
        (2, 2000, 8000, 1.05, '新秀', '#3498db'),
        (3, 5000, 15000, 1.05, '新秀', '#3498db'),
        (4, 10000, 30000, 1.1, '新秀', '#3498db'),
        (5, 20000, 50000, 1.15, '精英', '#9b59b6'),
        (10, 100000, 200000, 1.30, '大师', '#e74c3c'),
        (20, 500000, 1000000, 1.60, '传奇', '#f39c12'),
        (30, 1500000, 3000000, 2.00, '神话', '#e67e22')
    ON CONFLICT (level) DO NOTHING;

    INSERT INTO system_config (key, value, description) VALUES
        ('danmu_rate_limit', '20', '每分钟弹幕数限制'),
        ('gift_rate_limit', '30', '每分钟礼物数限制'),
        ('coin_recharge_min', '10', '最小充值金额'),
        ('stream_key_expire_days', '30', '推流密钥有效期（天）')
    ON CONFLICT (key) DO NOTHING;

    INSERT INTO sensitive_words (word, type, severity) VALUES
        ('赌博', 'blacklist', 'high'),
        ('色情', 'blacklist', 'high'),
        ('毒品', 'blacklist', 'high'),
        ('诈骗', 'blacklist', 'high')
    ON CONFLICT DO NOTHING;

    INSERT INTO gifts (name, coin_price, icon_url, animation_type, sort_order) VALUES
        ('鲜花', 10, '/gifts/flower.png', 'css', 1),
        ('爱心', 30, '/gifts/heart.png', 'css', 2),
        ('掌声', 50, '/gifts/clap.png', 'css', 3),
        ('火箭', 100, '/gifts/rocket.png', 'lottie', 4),
        ('游艇', 500, '/gifts/yacht.png', 'lottie', 5),
        ('飞机', 1000, '/gifts/plane.png', 'lottie', 6),
        ('钻戒', 5000, '/gifts/ring.png', 'particle', 7),
        ('城堡', 10000, '/gifts/castle.png', 'particle', 8)
    ON CONFLICT DO NOTHING;

    SELECT 'Database initialized successfully!' AS status;
EOSQL
