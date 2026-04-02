package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type TGConfig struct {
	Token      string
	Admins_Ids []int64
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Name     string
	Port     int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type YtdlpConfig struct {
	BinaryPath  string
	OutputDir   string
	MaxFileSize string
	Format      string
}

type Config struct {
	DB      DBConfig
	Redis   RedisConfig
	TG      TGConfig
	Ytdlp   YtdlpConfig
	Workers int
	Logger  string
}

func parseAdminIds(raw string) []int64 {
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		id, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			continue
		}

		ids = append(ids, id)
	}

	return ids
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, fmt.Errorf("Failed to load .env: %v", err)
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		dbPort = 5432
	}

	workersConfig, err := strconv.Atoi(os.Getenv("WORKERS_NUMBER"))
	if err != nil {
		workersConfig = 5
	}

	return &Config{
		DB: DBConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Name:     os.Getenv("DB_NAME"),
			Port:     dbPort,
		},
		Redis: RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDB,
		},
		TG: TGConfig{
			Token:      os.Getenv("TELEGRAM_TOKEN"),
			Admins_Ids: parseAdminIds(os.Getenv("TELEGRAM_ADMINS_IDS")),
		},
		Ytdlp: YtdlpConfig{
			BinaryPath:  os.Getenv("YTDLP_BINARYPATH"),
			OutputDir:   os.Getenv("YTDLP_OUTPUTDIR"),
			MaxFileSize: os.Getenv("YTDLP_MAXFILESIZE"),
			Format:      os.Getenv("YTDLP_FORMAT"),
		},
		Logger:  os.Getenv("LOGGER_LVL"),
		Workers: workersConfig,
	}, nil
}
