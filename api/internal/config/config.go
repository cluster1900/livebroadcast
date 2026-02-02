package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	AccessTTL  int
	RefreshTTL int
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

func Load() (*Config, error) {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8888"
	}
	mode := os.Getenv("SERVER_MODE")
	if mode == "" {
		mode = "debug"
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	if dbPort == 0 {
		dbPort = 5432
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "huya_live"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "huya_live_secret"
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "huya_live"
	}
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	if redisDB == 0 {
		redisDB = 0
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret_change_in_production"
	}
	accessTTL, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TTL"))
	if accessTTL == 0 {
		accessTTL = 900
	}
	refreshTTL, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TTL"))
	if refreshTTL == 0 {
		refreshTTL = 604800
	}

	return &Config{
		Server: ServerConfig{
			Port: port,
			Mode: mode,
		},
		Database: DatabaseConfig{
			Host:     host,
			Port:     dbPort,
			User:     user,
			Password: password,
			Name:     name,
			SSLMode:  sslMode,
		},
		Redis: RedisConfig{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       redisDB,
		},
		JWT: JWTConfig{
			Secret:     jwtSecret,
			AccessTTL:  accessTTL,
			RefreshTTL: refreshTTL,
		},
	}, nil
}
