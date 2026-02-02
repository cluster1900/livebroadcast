package main

import (
	"log"

	"github.com/huya_live/api/internal/config"
	"github.com/huya_live/api/internal/repository"
	"github.com/huya_live/api/internal/routes"
	"github.com/huya_live/api/pkg/redis"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化Redis
	if err := redis.Init(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB); err != nil {
		log.Fatalf("Failed to init Redis: %v", err)
	}

	// 初始化数据库
	if err := repository.InitDB(cfg.Database); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// 初始化Gin路由
	r := routes.SetupRouter()

	// 启动服务
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
