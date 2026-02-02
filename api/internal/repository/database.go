package repository

import (
	"fmt"
	"github.com/huya_live/api/internal/config"
	"github.com/huya_live/api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg config.DatabaseConfig) error {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := DB.AutoMigrate(
		&models.User{},
		&models.Streamer{},
		&models.LiveRoom{},
		&models.Gift{},
		&models.GiftTransaction{},
		&models.CoinTransaction{},
		&models.FanRelation{},
		&models.LevelConfig{},
		&models.SensitiveWord{},
		&models.SystemConfig{},
	); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := seedData(); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	return nil
}

func seedData() error {
	var count int64
	DB.Model(&models.LevelConfig{}).Count(&count)
	if count > 0 {
		return nil
	}

	levels := []models.LevelConfig{
		{Level: 1, ExpRequired: 0, LoyaltyPointsRequired: ptrInt64(0), BonusMultiplier: 1.0, LevelName: "萌新", Color: "#999999"},
		{Level: 2, ExpRequired: 2000, LoyaltyPointsRequired: ptrInt64(8000), BonusMultiplier: 1.05, LevelName: "新秀", Color: "#3498db"},
		{Level: 3, ExpRequired: 5000, LoyaltyPointsRequired: ptrInt64(15000), BonusMultiplier: 1.05, LevelName: "新秀", Color: "#3498db"},
		{Level: 4, ExpRequired: 10000, LoyaltyPointsRequired: ptrInt64(30000), BonusMultiplier: 1.1, LevelName: "新秀", Color: "#3498db"},
		{Level: 5, ExpRequired: 20000, LoyaltyPointsRequired: ptrInt64(50000), BonusMultiplier: 1.15, LevelName: "精英", Color: "#9b59b6"},
		{Level: 10, ExpRequired: 100000, LoyaltyPointsRequired: ptrInt64(200000), BonusMultiplier: 1.30, LevelName: "大师", Color: "#e74c3c"},
		{Level: 20, ExpRequired: 500000, LoyaltyPointsRequired: ptrInt64(1000000), BonusMultiplier: 1.60, LevelName: "传奇", Color: "#f39c12"},
		{Level: 30, ExpRequired: 1500000, LoyaltyPointsRequired: ptrInt64(3000000), BonusMultiplier: 2.00, LevelName: "神话", Color: "#e67e22"},
	}

	if err := DB.Create(&levels).Error; err != nil {
		return err
	}

	configs := []models.SystemConfig{
		{Key: "danmu_rate_limit", Value: "20", Description: "每分钟弹幕数限制"},
		{Key: "gift_rate_limit", Value: "30", Description: "每分钟礼物数限制"},
		{Key: "coin_recharge_min", Value: "10", Description: "最小充值金额"},
		{Key: "stream_key_expire_days", Value: "30", Description: "推流密钥有效期（天）"},
	}

	if err := DB.Create(&configs).Error; err != nil {
		return err
	}

	sensitiveWords := []models.SensitiveWord{
		{Word: "赌博", Type: "blacklist", Severity: "high"},
		{Word: "色情", Type: "blacklist", Severity: "high"},
		{Word: "毒品", Type: "blacklist", Severity: "high"},
		{Word: "诈骗", Type: "blacklist", Severity: "high"},
	}

	if err := DB.Create(&sensitiveWords).Error; err != nil {
		return err
	}

	gifts := []models.Gift{
		{Name: "鲜花", CoinPrice: 10, IconURL: "/gifts/flower.png", AnimationType: "css", SortOrder: 1},
		{Name: "爱心", CoinPrice: 30, IconURL: "/gifts/heart.png", AnimationType: "css", SortOrder: 2},
		{Name: "掌声", CoinPrice: 50, IconURL: "/gifts/clap.png", AnimationType: "css", SortOrder: 3},
		{Name: "火箭", CoinPrice: 100, IconURL: "/gifts/rocket.png", AnimationType: "lottie", SortOrder: 4},
		{Name: "游艇", CoinPrice: 500, IconURL: "/gifts/yacht.png", AnimationType: "lottie", SortOrder: 5},
		{Name: "飞机", CoinPrice: 1000, IconURL: "/gifts/plane.png", AnimationType: "lottie", SortOrder: 6},
		{Name: "钻戒", CoinPrice: 5000, IconURL: "/gifts/ring.png", AnimationType: "particle", SortOrder: 7},
		{Name: "城堡", CoinPrice: 10000, IconURL: "/gifts/castle.png", AnimationType: "particle", SortOrder: 8},
	}

	if err := DB.Create(&gifts).Error; err != nil {
		return err
	}

	return nil
}

func ptrInt64(v int64) *int64 {
	return &v
}
