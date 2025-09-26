package handler

import (
	"context"

	. "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
)

func StartCommand(ctx context.Context, bot *Bot, update *models.Update, client *ent.Client, log *zap.SugaredLogger) {
	log.Named("start").Info("hi")
	_, err := bot.SendMessage(ctx, &SendMessageParams{ChatID: update.Message.From.ID, Text: "Hi"})
	if err != nil {
		panic("")
	}
}
