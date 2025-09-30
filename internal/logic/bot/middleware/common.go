package middleware

import (
	"github.com/go-telegram/bot"
)

func Middlewares() []bot.Middleware {
	return []bot.Middleware{
		loggingMiddleware,
		recoverMiddleware,
	}
}
