package bot

import (
	"context"
	"run_ruby_bot/core"
	"run_ruby_bot/service"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func startHandler(ctx context.Context, b *tg.Bot, update *models.Update) {
	b.SendMessage(ctx, &tg.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hello! I can interpret Ruby codes.",
	})
}

func rubyCodeHandler(ctx context.Context, b *tg.Bot, update *models.Update) {
	userId := update.Message.From.ID
	if !core.Cfg().Bot.AllowedUsersSet[userId] {
		b.SendMessage(ctx, &tg.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "<em>permission denied</em>",
			ParseMode: models.ParseModeHTML,
		})
		return
	}

	msg, err := b.SendMessage(ctx, &tg.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "<em>waiting...</em>",
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		return
	}

	code := service.SourceCode(update.Message.Text)
	ch := make(chan service.HTMLMessage)
	go service.RunInterpretTask(code, ch)

	for text := range ch {
		b.EditMessageText(ctx, &tg.EditMessageTextParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: msg.ID,
			Text:      string(text),
			ParseMode: models.ParseModeHTML,
		})
	}
}
