package bot

import (
	"log/slog"
	"run_ruby_bot/core"

	tg "github.com/go-telegram/bot"
)

func NewRunRubyBot() *tg.Bot {
	opts := []tg.Option{
		tg.WithDefaultHandler(rubyCodeHandler),
	}

	b, err := tg.New(core.Cfg().Bot.Token, opts...)
	if err != nil {
		slog.Error("failed to create the telegram bot,", slog.String("error", err.Error()))
		panic("initialization failed")
	}
	b.RegisterHandler(tg.HandlerTypeMessageText, "/start", tg.MatchTypeExact, startHandler)
	return b
}
