package bot

import (
	"github.com/go-telegram/bot"
	"go.uber.org/fx"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/bot/handler"
)

func Handlers() fx.Option {
	return fx.Provide(
		AsHandler(NewTelegramHandler(StartCommand, &TelegramHandlerOpts{
			Handler: bot.HandlerTypeMessageText,
			Pattern: "/start",
			Match:   bot.MatchTypeExact,
		})),
	)
}
