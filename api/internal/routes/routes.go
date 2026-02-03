package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/huya_live/api/internal/handlers"
	"github.com/huya_live/api/internal/middleware"
	"github.com/huya_live/api/pkg/centrifugo"
	"github.com/huya_live/api/pkg/jwt"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Next()
	})

	jwtManager := jwt.NewManager(
		"your_super_secret_key_change_in_production",
		900,
		604800,
	)

	centrifugoClient := centrifugo.NewClient("http://localhost:8000", "your_centrifugo_api_key")

	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(jwtManager)
	liveHandler := handlers.NewLiveHandler()
	streamerHandler := handlers.NewStreamerHandler()
	srsHandler := handlers.NewSRSHandler()
	centrifugoHandler := handlers.NewCentrifugoHandler()
	danmuHandler := handlers.NewDanmuHandler()
	giftHandler := handlers.NewGiftHandler()
	walletHandler := handlers.NewWalletHandler()
	socialHandler := handlers.NewSocialHandler()
	relayHandler := handlers.NewRelayHandler()
	tvHandler := handlers.NewPredefinedTVHandler()
	leaderboardHandler := handlers.NewLeaderboardHandler()
	notificationHandler := handlers.NewNotificationHandler()
	messageHandler := handlers.NewMessageHandler(centrifugoClient)
	historyHandler := handlers.NewHistoryHandler()
	reportHandler := handlers.NewReportHandler()
	giftInventoryHandler := handlers.NewGiftInventoryHandler()
	likeHandler := handlers.NewLikeHandler()
	scheduleHandler := handlers.NewScheduleHandler()
	passwordHandler := handlers.NewPasswordHandler()
	adminHandler := handlers.NewAdminHandler()

	r.GET("/health", healthHandler.HealthCheck)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		users := api.Group("/user")
		{
			users.GET("/profile", middleware.JWTRequired(jwtManager), authHandler.GetProfile)
			users.PUT("/profile", middleware.JWTRequired(jwtManager), authHandler.UpdateProfile)
		}

		live := api.Group("/live")
		{
			live.GET("/rooms", liveHandler.ListRooms)
			live.GET("/rooms/:id", liveHandler.GetRoom)
			live.GET("/live-rooms", liveHandler.ListRooms)
		}

		streamers := api.Group("/streamers")
		{
			streamers.POST("/apply", middleware.JWTRequired(jwtManager), streamerHandler.Apply)
			streamers.GET("/me", middleware.JWTRequired(jwtManager), streamerHandler.GetInfo)
			streamers.POST("/refresh-key", middleware.JWTRequired(jwtManager), streamerHandler.RefreshStreamKey)
		}

		rooms := api.Group("/rooms")
		rooms.Use(middleware.JWTRequired(jwtManager))
		{
			rooms.POST("", liveHandler.CreateRoom)
			rooms.PUT("/:id", liveHandler.UpdateRoom)
			rooms.POST("/:id/end", liveHandler.EndRoom)
		}

		centrifugo := api.Group("/centrifugo")
		{
			centrifugo.GET("/token", middleware.JWTRequired(jwtManager), centrifugoHandler.GetToken)
		}

		danmu := api.Group("/danmu")
		danmu.Use(middleware.JWTRequired(jwtManager))
		{
			danmu.POST("/send", danmuHandler.SendDanmu)
		}

		gifts := api.Group("/gifts")
		gifts.Use(middleware.JWTRequired(jwtManager))
		{
			gifts.POST("/send", giftHandler.SendGift)
		}

		wallet := api.Group("/wallet")
		wallet.Use(middleware.JWTRequired(jwtManager))
		{
			wallet.POST("/recharge", walletHandler.RechargeCoins)
			wallet.GET("/balance", walletHandler.GetBalance)
			wallet.GET("/transactions", walletHandler.GetTransactionHistory)
		}

		social := api.Group("/social")
		{
			social.POST("/follow", middleware.JWTRequired(jwtManager), socialHandler.Follow)
			social.POST("/unfollow", middleware.JWTRequired(jwtManager), socialHandler.Unfollow)
			social.GET("/followings", middleware.JWTRequired(jwtManager), socialHandler.GetFollowings)
			social.GET("/followers/:streamer_id", socialHandler.GetFollowers)
		}

		relay := api.Group("/relay")
		{
			relay.GET("", relayHandler.GetRelays)
			relay.GET("/:id", relayHandler.GetRelay)
			relay.POST("", relayHandler.CreateRelay)
			relay.PUT("/:id", relayHandler.UpdateRelay)
			relay.DELETE("/:id", relayHandler.DeleteRelay)
			relay.POST("/:id/start", relayHandler.StartRelay)
			relay.POST("/:id/stop", relayHandler.StopRelay)
		}

		tv := api.Group("/tv")
		{
			tv.GET("", tvHandler.GetTVStations)
			tv.POST("", tvHandler.AddTVStation)
			tv.POST("/import", tvHandler.CreateRelaysFromTVStations)
		}

		leaderboard := api.Group("/leaderboard")
		{
			leaderboard.GET("/rooms/:room_id", leaderboardHandler.GetRoomLeaderboard)
			leaderboard.GET("/global", leaderboardHandler.GetGlobalLeaderboard)
			leaderboard.GET("/rich", leaderboardHandler.GetRichList)
		}

		extra := api.Group("/extra")
		{
			extra.GET("/categories", leaderboardHandler.GetCategories)
			extra.GET("/online-count", leaderboardHandler.GetOnlineCount)
		}

		notifications := api.Group("/notifications")
		notifications.Use(middleware.JWTRequired(jwtManager))
		{
			notifications.GET("", notificationHandler.GetNotifications)
			notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
			notifications.POST("/:id/read", notificationHandler.MarkAsRead)
			notifications.POST("/read-all", notificationHandler.MarkAllAsRead)
			notifications.DELETE("/:id", notificationHandler.DeleteNotification)
		}

		messages := api.Group("/messages")
		messages.Use(middleware.JWTRequired(jwtManager))
		{
			messages.POST("/send", messageHandler.SendMessage)
			messages.GET("/conversations", messageHandler.GetConversations)
			messages.GET("/with/:user_id", messageHandler.GetMessages)
			messages.GET("/unread-count", messageHandler.GetUnreadMessageCount)
			messages.DELETE("/conversation/:user_id", messageHandler.DeleteConversation)
		}

		history := api.Group("/history")
		history.Use(middleware.JWTRequired(jwtManager))
		{
			history.GET("/watch", historyHandler.GetWatchHistory)
			history.POST("/watch", historyHandler.AddWatchHistory)
			history.DELETE("/watch", historyHandler.ClearWatchHistory)
			history.DELETE("/watch/:id", historyHandler.DeleteWatchHistory)
		}

		reports := api.Group("/reports")
		reports.Use(middleware.JWTRequired(jwtManager))
		{
			reports.POST("", reportHandler.CreateReport)
			reports.GET("/my", reportHandler.GetMyReports)
		}

		reportsAdmin := api.Group("/admin/reports")
		reportsAdmin.Use(middleware.JWTRequired(jwtManager))
		{
			reportsAdmin.GET("/pending", reportHandler.GetPendingReports)
			reportsAdmin.POST("/:id/handle", reportHandler.HandleReport)
		}

		inventory := api.Group("/inventory")
		inventory.Use(middleware.JWTRequired(jwtManager))
		{
			inventory.GET("/gifts", giftInventoryHandler.GetInventory)
			inventory.POST("/use", giftInventoryHandler.UseGift)
		}

		likes := api.Group("/likes")
		{
			likes.POST("/rooms/:room_id", middleware.JWTRequired(jwtManager), likeHandler.LikeRoom)
			likes.DELETE("/rooms/:room_id", middleware.JWTRequired(jwtManager), likeHandler.UnlikeRoom)
			likes.GET("/rooms/:room_id/count", likeHandler.GetLikeCount)
			likes.GET("/rooms/:room_id/status", middleware.JWTRequired(jwtManager), likeHandler.HasLiked)
		}

		schedules := api.Group("/schedules")
		schedules.Use(middleware.JWTRequired(jwtManager))
		{
			schedules.POST("", scheduleHandler.CreateSchedule)
			schedules.GET("/my", scheduleHandler.GetMySchedules)
			schedules.PUT("/:id", scheduleHandler.UpdateSchedule)
			schedules.POST("/:id/cancel", scheduleHandler.CancelSchedule)
			schedules.DELETE("/:id", scheduleHandler.DeleteSchedule)
		}

		extraSchedules := api.Group("/extra/schedules")
		{
			extraSchedules.GET("/upcoming", scheduleHandler.GetUpcomingSchedules)
		}

		password := api.Group("/password")
		{
			password.POST("/change", middleware.JWTRequired(jwtManager), passwordHandler.ChangePassword)
			password.POST("/reset/request", passwordHandler.RequestReset)
			password.POST("/reset/complete", passwordHandler.CompleteReset)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.JWTRequired(jwtManager))
		{
			admin.GET("/dashboard", adminHandler.GetDashboardStats)
			admin.GET("/users", adminHandler.GetUserList)
			admin.POST("/users/:id/ban", adminHandler.BanUser)
			admin.POST("/users/:id/unban", adminHandler.UnbanUser)
			admin.GET("/rooms", adminHandler.GetRoomList)
			admin.POST("/rooms/:id/ban", adminHandler.BanRoom)
			admin.GET("/gifts", adminHandler.GetGiftList)
			admin.POST("/gifts", adminHandler.CreateGift)
			admin.PUT("/gifts/:id", adminHandler.UpdateGift)
			admin.DELETE("/gifts/:id", adminHandler.DeleteGift)
			admin.GET("/sensitive-words", adminHandler.GetSensitiveWords)
			admin.POST("/sensitive-words", adminHandler.AddSensitiveWord)
			admin.DELETE("/sensitive-words/:id", adminHandler.DeleteSensitiveWord)
			admin.GET("/config", adminHandler.GetSystemConfig)
			admin.PUT("/config", adminHandler.UpdateSystemConfig)
		}
	}

	srs := r.Group("/api/srs")
	{
		srs.POST("/callback/publish", srsHandler.OnPublish)
		srs.POST("/callback/unpublish", srsHandler.OnUnpublish)
	}

	return r
}
