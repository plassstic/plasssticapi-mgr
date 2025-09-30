package middleware

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/depend/logger"
)

func recoverMiddleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		defer func() {
			if r := recover(); r != nil {
				GetLogger("bot.middleware.recover").Errorf("recovered from panic while updating; (%e)", r)
			}
		}()

		next(ctx, b, update)
	}
}
