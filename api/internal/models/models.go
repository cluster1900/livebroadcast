package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username     string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	Nickname     string     `gorm:"type:varchar(100)" json:"nickname"`
	AvatarURL    string     `gorm:"type:text" json:"avatar_url"`
	Phone        string     `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
	Email        string     `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Level        int        `gorm:"default:1" json:"level"`
	Exp          int64      `gorm:"default:0" json:"exp"`
	CoinBalance  int        `gorm:"default:0" json:"coin_balance"`
	Status       string     `gorm:"type:varchar(20);default:'active'" json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type Streamer struct {
	UserID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id"`
	StreamKey         string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"stream_key"`
	StreamKeyExpireAt *time.Time `json:"stream_key_expire_at"`
	RtmpURL           string     `gorm:"type:text" json:"rtmp_url"`
	Status            string     `gorm:"type:varchar(20);default:'offline'" json:"status"`
	IsVerified        bool       `gorm:"default:false" json:"is_verified"`
	TotalRevenue      int64      `gorm:"default:0" json:"total_revenue"`
	FollowerCount     int        `gorm:"default:0" json:"follower_count"`
	TotalLiveDuration int        `gorm:"default:0" json:"total_live_duration"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

type LiveRoom struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StreamerID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"streamer_id"`
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Category    string     `gorm:"type:varchar(50)" json:"category"`
	CoverURL    string     `gorm:"type:text" json:"cover_url"`
	ChannelName string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"channel_name"`
	Status      string     `gorm:"type:varchar(20);default:'ended'" json:"status"`
	StartAt     *time.Time `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	PeakOnline  int        `gorm:"default:0" json:"peak_online"`
	TotalViews  int        `gorm:"default:0" json:"total_views"`
	RecordURL   string     `gorm:"type:text" json:"record_url"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

type Gift struct {
	ID               int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string    `gorm:"type:varchar(50);not null" json:"name"`
	CoinPrice        int       `gorm:"not null" json:"coin_price"`
	IconURL          string    `gorm:"type:text;not null" json:"icon_url"`
	AnimationType    string    `gorm:"type:varchar(20)" json:"animation_type"`
	AnimationURL     string    `gorm:"type:text" json:"animation_url"`
	MinLevelRequired int       `gorm:"default:1" json:"min_level_required"`
	IsActive         bool      `gorm:"default:true" json:"is_active"`
	SortOrder        int       `gorm:"default:0" json:"sort_order"`
	Category         string    `gorm:"type:varchar(20);default:'normal'" json:"category"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type GiftTransaction struct {
	ID                  int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID            uuid.UUID `gorm:"type:uuid;not null;index" json:"sender_id"`
	ReceiverID          uuid.UUID `gorm:"type:uuid;not null;index" json:"receiver_id"`
	RoomID              uuid.UUID `gorm:"type:uuid;not null;index" json:"room_id"`
	GiftID              int       `gorm:"not null;index" json:"gift_id"`
	GiftCount           int       `gorm:"default:1" json:"gift_count"`
	CoinAmount          int       `gorm:"not null" json:"coin_amount"`
	LoyaltyPointsGained int64     `json:"loyalty_points_gained"`
	UserLevelAtSend     int       `json:"user_level_at_send"`
	BonusMultiplier     float64   `gorm:"default:1.0" json:"bonus_multiplier"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type CoinTransaction struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Amount       int       `gorm:"not null" json:"amount"`
	BalanceAfter int       `gorm:"not null" json:"balance_after"`
	Type         string    `gorm:"type:varchar(20);not null;index" json:"type"`
	RelatedID    *int64    `json:"related_id"`
	Description  string    `gorm:"type:text" json:"description"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type FanRelation struct {
	UserID          uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	StreamerID      uuid.UUID  `gorm:"type:uuid;not null" json:"streamer_id"`
	FanLevel        int        `gorm:"default:1" json:"fan_level"`
	LoyaltyPoints   int64      `gorm:"default:0" json:"loyalty_points"`
	BadgeName       string     `gorm:"type:varchar(20)" json:"badge_name"`
	BadgeWorn       bool       `gorm:"default:true" json:"badge_worn"`
	TotalGiftAmount int64      `gorm:"default:0" json:"total_gift_amount"`
	FollowedAt      time.Time  `gorm:"autoCreateTime" json:"followed_at"`
	LastGiftAt      *time.Time `json:"last_gift_at"`
}

type LevelConfig struct {
	Level                 int     `gorm:"primaryKey" json:"level"`
	ExpRequired           int64   `gorm:"not null" json:"exp_required"`
	LoyaltyPointsRequired *int64  `json:"loyalty_points_required"`
	BonusMultiplier       float64 `gorm:"default:1.0" json:"bonus_multiplier"`
	LevelName             string  `gorm:"type:varchar(50)" json:"level_name"`
	IconURL               string  `gorm:"type:text" json:"icon_url"`
	Color                 string  `gorm:"type:varchar(7)" json:"color"`
}

type SensitiveWord struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Word      string    `gorm:"type:varchar(50);not null;index" json:"word"`
	Type      string    `gorm:"type:varchar(20);default:'blacklist'" json:"type"`
	Severity  string    `gorm:"type:varchar(10);default:'medium'" json:"severity"`
	IsActive  bool      `gorm:"default:true;index" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type SystemConfig struct {
	Key         string    `gorm:"type:varchar(50);primaryKey" json:"key"`
	Value       string    `gorm:"type:text;not null" json:"value"`
	Description string    `gorm:"type:text" json:"description"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
