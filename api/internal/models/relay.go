package models

import (
	"time"

	"github.com/google/uuid"
)

type RelayStream struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	Description   string    `gorm:"type:text" json:"description"`
	SourceURL     string    `gorm:"type:text;not null" json:"source_url"`
	SourceType    string    `gorm:"type:varchar(20);not null;default:'rtmp'" json:"source_type"`
	ChannelName   string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"channel_name"`
	StreamKey     string    `gorm:"type:varchar(128);not null" json:"stream_key"`
	RelayProtocol string    `gorm:"type:varchar(20);default:'rtmp'" json:"relay_protocol"`
	Status        string    `gorm:"type:varchar(20);default:'stopped'" json:"status"`
	Category      string    `gorm:"type:varchar(50)" json:"category"`
	CoverURL      string    `gorm:"type:text" json:"cover_url"`
	ViewCount     int64     `gorm:"default:0" json:"view_count"`
	PeakOnline    int64     `gorm:"default:0" json:"peak_online"`
	IsVerified    bool      `gorm:"default:false" json:"is_verified"`
	AutoStart     bool      `gorm:"default:true" json:"auto_start"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type RelayStreamLog struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	RelayStreamID uuid.UUID `gorm:"type:uuid;index;not null" json:"relay_stream_id"`
	EventType     string    `gorm:"type:varchar(50)" json:"event_type"`
	EventData     string    `gorm:"type:text" json:"event_data"`
	ErrorMessage  string    `gorm:"type:text" json:"error_message"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type PredefinedRelay struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	SourceURL   string    `gorm:"type:text;not null" json:"source_url"`
	Category    string    `gorm:"type:varchar(50)" json:"category"`
	Country     string    `gorm:"type:varchar(50)" json:"country"`
	Language    string    `gorm:"type:varchar(50)" json:"language"`
	CoverURL    string    `gorm:"type:text" json:"cover_url"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
