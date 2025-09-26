package handler

import (
	"context"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	. "plassstic.tech/gopkg/golang-manager/internal/depend/logger"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
)

func StartCommand(ctx context.Context, bot *tg.Bot, update *models.Update, client *ent.Client) {
	log := GetLogger("bot.handlers.StartCommand")

	log.Info("catch")

	_, err := bot.SendMessage(ctx, &tg.SendMessageParams{ChatID: update.Message.From.ID, Text: "Hi"})
	if err != nil {
		panic("")
	}
}
