package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"run_ruby_bot/bot"
	"run_ruby_bot/core"
	"run_ruby_bot/service"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("failed to load .env file")
	}

	core.InitCfg()
	core.InitInterpreter()
	service.InitService()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	b := bot.NewRunRubyBot()

	slog.Info("running bot...")
	b.Start(ctx)
}
