package bot

import (
	"log/slog"
	"run_ruby_bot/core"

	tg "github.com/go-telegram/bot"
)

func NewRunRubyBot() *tg.Bot {
	opts := []tg.Option{
		tg.WithDefaultHandler(DefaultHandler),
	}

	b, err := tg.New(core.Cfg().Bot.Token, opts...)
	if err != nil {
		slog.Error("failed to create the telegram bot,", "error", err.Error())
		panic("initialization failed")
	}
	b.RegisterHandler(tg.HandlerTypeMessageText, "/start", tg.MatchTypeExact, CmdStartHandler)
	b.RegisterHandler(tg.HandlerTypeMessageText, "/ruby ", tg.MatchTypePrefix, CmdRubyHandler)
	b.RegisterHandler(tg.HandlerTypeMessageText, "/ruby\n", tg.MatchTypePrefix, CmdRubyHandler)
	return b
}
