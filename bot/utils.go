package bot

import (
	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func newReply(to *models.Message, text string) *tg.SendMessageParams {
	return &tg.SendMessageParams{
		ReplyParameters: &models.ReplyParameters{
			MessageID:                to.ID,
			ChatID:                   to.Chat.ID,
			AllowSendingWithoutReply: true,
		},
		ChatID:    to.Chat.ID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	}
}
