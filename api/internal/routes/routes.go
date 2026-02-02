package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/huya_live/api/internal/handlers"
	"github.com/huya_live/api/internal/middleware"
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
	}

	srs := r.Group("/api/srs")
	{
		srs.POST("/callback/publish", srsHandler.OnPublish)
		srs.POST("/callback/unpublish", srsHandler.OnUnpublish)
	}

	return r
}
