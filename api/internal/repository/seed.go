package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/huya_live/api/internal/models"
)

func SeedTestData() error {
	var count int64

	count = 0
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	users := []models.User{
		{
			ID:           uuid.New(),
			Username:     "testuser1",
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.NT5.0op.W1yHqXSQae",
			Nickname:     "测试用户1",
			AvatarURL:    "https://api.dicebear.com/7.x/avataaars/svg?seed=testuser1",
			Email:        "test1@example.com",
			Level:        5,
			Exp:          25000,
			CoinBalance:  5000,
			Status:       "active",
		},
		{
			ID:           uuid.New(),
			Username:     "testuser2",
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.NT5.0op.W1yHqXSQae",
			Nickname:     "测试主播",
			AvatarURL:    "https://api.dicebear.com/7.x/avataaars/svg?seed=testuser2",
			Email:        "test2@example.com",
			Level:        10,
			Exp:          120000,
			CoinBalance:  10000,
			Status:       "active",
		},
		{
			ID:           uuid.New(),
			Username:     "testuser3",
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.NT5.0op.W1yHqXSQae",
			Nickname:     "主播小姐姐",
			AvatarURL:    "https://api.dicebear.com/7.x/avataaars/svg?seed=testuser3",
			Email:        "test3@example.com",
			Level:        8,
			Exp:          80000,
			CoinBalance:  8000,
			Status:       "active",
		},
		{
			ID:           uuid.New(),
			Username:     "admin",
			PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.NT5.0op.W1yHqXSQae",
			Nickname:     "管理员",
			AvatarURL:    "https://api.dicebear.com/7.x/avataaars/svg?seed=admin",
			Email:        "admin@example.com",
			Level:        30,
			Exp:          2000000,
			CoinBalance:  50000,
			Status:       "active",
		},
	}

	for i := range users {
		users[i].CreatedAt = time.Now()
		users[i].UpdatedAt = time.Now()
	}

	if err := DB.Create(&users).Error; err != nil {
		return err
	}

	streamers := []models.Streamer{
		{UserID: users[1].ID, StreamKey: "sk_live_" + uuid.New().String()[:8], RtmpURL: "rtmp://localhost/live", Status: "offline", IsVerified: true, FollowerCount: 100},
		{UserID: users[2].ID, StreamKey: "sk_live_" + uuid.New().String()[:8], RtmpURL: "rtmp://localhost/live", Status: "offline", IsVerified: true, FollowerCount: 50},
	}
	for i := range streamers {
		streamers[i].CreatedAt = time.Now()
	}
	if err := DB.Create(&streamers).Error; err != nil {
		return err
	}

	liveRooms := []models.LiveRoom{
		{StreamerID: users[1].ID, Title: "游戏直播：英雄联盟高端局", Category: "游戏", ChannelName: "room_" + uuid.New().String()[:8], Status: "live", StartAt: ptrTime(time.Now()), PeakOnline: 500, TotalViews: 10000},
		{StreamerID: users[2].ID, Title: "唱歌直播：流行金曲", Category: "音乐", ChannelName: "room_" + uuid.New().String()[:8], Status: "live", StartAt: ptrTime(time.Now()), PeakOnline: 300, TotalViews: 5000},
		{StreamerID: users[1].ID, Title: "户外直播：城市探险", Category: "户外", ChannelName: "room_" + uuid.New().String()[:8], Status: "ended", PeakOnline: 200, TotalViews: 3000},
	}
	for i := range liveRooms {
		liveRooms[i].CreatedAt = time.Now()
	}
	if err := DB.Create(&liveRooms).Error; err != nil {
		return err
	}

	gifts := []models.Gift{
		{Name: "鲜花", CoinPrice: 10, IconURL: "🌸", AnimationType: "css", SortOrder: 1, Category: "normal"},
		{Name: "爱心", CoinPrice: 30, IconURL: "❤️", AnimationType: "css", SortOrder: 2, Category: "normal"},
		{Name: "掌声", CoinPrice: 50, IconURL: "👏", AnimationType: "css", SortOrder: 3, Category: "normal"},
		{Name: "火箭", CoinPrice: 100, IconURL: "🚀", AnimationType: "lottie", SortOrder: 4, Category: "normal"},
		{Name: "游艇", CoinPrice: 500, IconURL: "🛥️", AnimationType: "lottie", SortOrder: 5, Category: "vip"},
		{Name: "飞机", CoinPrice: 1000, IconURL: "✈️", AnimationType: "lottie", SortOrder: 6, Category: "vip"},
		{Name: "钻戒", CoinPrice: 5000, IconURL: "💍", AnimationType: "particle", SortOrder: 7, Category: "special"},
		{Name: "城堡", CoinPrice: 10000, IconURL: "🏰", AnimationType: "particle", SortOrder: 8, Category: "special"},
		{Name: "天使翅膀", CoinPrice: 200, IconURL: "🪽", AnimationType: "lottie", SortOrder: 9, Category: "normal"},
		{Name: "豪华跑车", CoinPrice: 2000, IconURL: "🏎️", AnimationType: "lottie", SortOrder: 10, Category: "vip"},
	}
	for i := range gifts {
		gifts[i].IsActive = true
		gifts[i].MinLevelRequired = 1
	}
	if err := DB.Create(&gifts).Error; err != nil {
		return err
	}

	fanRelations := []models.FanRelation{
		{UserID: users[0].ID, StreamerID: users[1].ID, FanLevel: 3, LoyaltyPoints: 50000, BadgeName: "铁粉"},
		{UserID: users[0].ID, StreamerID: users[2].ID, FanLevel: 1, LoyaltyPoints: 5000, BadgeName: "新粉"},
		{UserID: users[3].ID, StreamerID: users[1].ID, FanLevel: 5, LoyaltyPoints: 200000, BadgeName: "铁粉"},
	}
	for i := range fanRelations {
		fanRelations[i].FollowedAt = time.Now()
	}
	if err := DB.Create(&fanRelations).Error; err != nil {
		return err
	}

	notifications := []models.Notification{
		{UserID: users[0].ID, Type: "system", Title: "欢迎来到虎牙直播", Content: "感谢您注册虎牙直播，祝您直播愉快！", Link: "/"},
		{UserID: users[0].ID, Type: "gift", Title: "礼物到账", Content: "您收到了一份礼物！", Link: "/inventory"},
		{UserID: users[0].ID, Type: "follow", Title: "新粉丝", Content: "主播小姐姐成为了您的新粉丝！", Link: "/profile/2"},
	}
	for i := range notifications {
		notifications[i].IsRead = false
		notifications[i].CreatedAt = time.Now().Add(-time.Hour * time.Duration(i+1))
	}
	if err := DB.Create(&notifications).Error; err != nil {
		return err
	}

	privateMessages := []models.PrivateMessage{
		{SenderID: users[1].ID, ReceiverID: users[0].ID, Content: "感谢关注我的直播间！", IsRead: false},
		{SenderID: users[0].ID, ReceiverID: users[1].ID, Content: "主播加油！", IsRead: true},
		{SenderID: users[2].ID, ReceiverID: users[0].ID, Content: "欢迎来我的直播间玩~", IsRead: false},
	}
	for i := range privateMessages {
		privateMessages[i].CreatedAt = time.Now().Add(-time.Hour * time.Duration(i+1))
	}
	if err := DB.Create(&privateMessages).Error; err != nil {
		return err
	}

	watchHistories := []models.WatchHistory{
		{UserID: users[0].ID, RoomID: liveRooms[0].ID, WatchDuration: 3600},
		{UserID: users[0].ID, RoomID: liveRooms[1].ID, WatchDuration: 1800},
	}
	for i := range watchHistories {
		watchHistories[i].CreatedAt = time.Now().Add(-time.Hour * 24)
	}
	if err := DB.Create(&watchHistories).Error; err != nil {
		return err
	}

	liveSchedules := []models.LiveSchedule{
		{StreamerID: users[1].ID, Title: "今晚8点：精彩游戏直播", Description: "不见不散！", Category: "游戏", StartTime: time.Now().Add(time.Hour * 24), Status: "scheduled"},
		{StreamerID: users[2].ID, Title: "周末特别节目", Description: "准备了神秘惊喜", Category: "娱乐", StartTime: time.Now().Add(time.Hour * 48), Status: "scheduled"},
	}
	if err := DB.Create(&liveSchedules).Error; err != nil {
		return err
	}

	sensitiveWords := []models.SensitiveWord{
		{Word: "垃圾广告", Type: "blacklist", Severity: "high"},
		{Word: "恶意灌水", Type: "blacklist", Severity: "medium"},
		{Word: "违规内容", Type: "blacklist", Severity: "high"},
		{Word: "敏感词", Type: "blacklist", Severity: "low"},
	}
	if err := DB.Create(&sensitiveWords).Error; err != nil {
		return err
	}

	giftTransactions := []models.GiftTransaction{
		{SenderID: users[0].ID, ReceiverID: users[1].ID, RoomID: liveRooms[0].ID, GiftID: 4, GiftCount: 5, CoinAmount: 500, LoyaltyPointsGained: 500},
		{SenderID: users[0].ID, ReceiverID: users[2].ID, RoomID: liveRooms[1].ID, GiftID: 2, GiftCount: 10, CoinAmount: 300, LoyaltyPointsGained: 300},
	}
	for i := range giftTransactions {
		giftTransactions[i].UserLevelAtSend = 5
		giftTransactions[i].BonusMultiplier = 1.15
		giftTransactions[i].CreatedAt = time.Now().Add(-time.Hour * 12)
	}
	if err := DB.Create(&giftTransactions).Error; err != nil {
		return err
	}

	coinTransactions := []models.CoinTransaction{
		{UserID: users[0].ID, Amount: 1000, BalanceAfter: 5000, Type: "recharge", Description: "充值虎牙币"},
		{UserID: users[0].ID, Amount: -500, BalanceAfter: 4500, Type: "gift", Description: "赠送礼物"},
		{UserID: users[0].ID, Amount: 2000, BalanceAfter: 6500, Type: "recharge", Description: "充值虎牙币"},
	}
	for i := range coinTransactions {
		coinTransactions[i].CreatedAt = time.Now().Add(-time.Hour * 24 * time.Duration(i+1))
	}
	if err := DB.Create(&coinTransactions).Error; err != nil {
		return err
	}

	return nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func nullString(s string) *string {
	return &s
}
