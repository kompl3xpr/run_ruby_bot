package core

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Bot struct {
		Token string
		AllowedUsers []int64
		AllowedUsersSet map[int64]bool
		MessageMaxLength int
	}
	Task struct {
		Timeout time.Duration
		TaskPoolCapacity int
		QueueCapacity int
	}
	Docker struct {
		Name   string
		Memory string
		Cpus   float64
	}
}

var instance *Config

func InitCfg() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	slog.Info("reading configurations from file...")
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("failed to read configurations, ", slog.String("error", err.Error()))
		os.Exit(1)
	}



	var config Config
	slog.Info("parsing configurations...")
	if err := viper.Unmarshal(&config); err != nil {
		slog.Error("failed to parse configurations, ", slog.String("error", err.Error()))
		os.Exit(2)
	}

	allowedUsersSet := make(map[int64]bool)
	for _, user := range config.Bot.AllowedUsers {
		allowedUsersSet[int64(user)] = true
	}
	config.Bot.AllowedUsersSet = allowedUsersSet

	if len(config.Bot.Token) == 0 {
		slog.Warn("using environment variable `TELEGRAM_BOT_TOKEN`...")
		token := os.Getenv("TELEGRAM_BOT_TOKEN")
		if len(token) == 0 {
			slog.Error("environment variable `TELEGRAM_BOT_TOKEN` not found")
			os.Exit(3)
		}
		config.Bot.Token = token
	}

	instance = &config
	slog.Info("", slog.String("configuration", fmt.Sprintf("%#v", config)))
}

func Cfg() *Config {
	return instance
}
