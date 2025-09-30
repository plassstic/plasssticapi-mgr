package middleware

import (
	"context"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/depend/logger"
)

func loggingMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		t := time.Now().UnixMilli()
		next(ctx, b, update)
		elapsed := time.Now().UnixMilli() - t

		log := GetLogger("bot.middleware.logging")

		var userID int64

		userID = -1
		if update.Message != nil {
			if update.Message.From != nil {
				userID = update.Message.From.ID
			}
		}

		if update.CallbackQuery != nil {
			userID = update.CallbackQuery.From.ID
		}

		log.Infof("handled update %d by user %v in %d ms", update.ID, userID, elapsed)
	}
}
