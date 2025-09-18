package bot

import (
	"context"
	"run_ruby_bot/core"
	"run_ruby_bot/service"
	"strings"
	"unicode"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func CmdStartHandler(ctx context.Context, b *tg.Bot, update *models.Update) {
	if update.Message.Chat.Type != "private" {
		return
	}

	b.SendMessage(ctx, &tg.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hello! I can interpret Ruby codes.",
	})
}

func CmdRubyHandler(ctx context.Context, b *tg.Bot, update *models.Update) {
	runRubyCodeHandler(ctx, b, update, false, true)
}

func DefaultHandler(ctx context.Context, b *tg.Bot, update *models.Update) {
	runRubyCodeHandler(ctx, b, update, true, false)
}

func runRubyCodeHandler(
	ctx context.Context,
	b *tg.Bot,
	update *models.Update,
	checkPrivate bool,
	removeCmdPrefix bool,
) {
	var msg_in *models.Message
	if update.Message != nil {
		msg_in = update.Message
	} else if update.EditedMessage != nil {
		msg_in = update.EditedMessage
	} else {
		return
	}

	if checkPrivate && msg_in.Chat.Type != "private" {
		return
	}

	code := msg_in.Text
	if removeCmdPrefix {
		code = strings.TrimLeftFunc(code, func(c rune) bool { return !unicode.IsSpace(c) })
		code = strings.TrimLeftFunc(code, unicode.IsSpace)
	}
	if len(strings.TrimSpace(code)) == 0 {
		return
	}

	userId := msg_in.From.ID
	if core.Cfg().Bot.Whitelist && !core.Cfg().Bot.AllowedUsersSet[userId] {
		b.SendMessage(ctx, newReply(msg_in, "<em>permission denied</em>"))
		return
	}

	msg_out, err := b.SendMessage(ctx, newReply(msg_in, "<em>waiting...</em>"))
	if err != nil {
		return
	}

	ch := make(chan service.HTMLMessage)
	go service.RunInterpretTask(service.SourceCode(code), ch)

	for text := range ch {
		b.EditMessageText(ctx, &tg.EditMessageTextParams{
			ChatID:    msg_in.Chat.ID,
			MessageID: msg_out.ID,
			Text:      string(text),
			ParseMode: models.ParseModeHTML,
		})
	}
}
