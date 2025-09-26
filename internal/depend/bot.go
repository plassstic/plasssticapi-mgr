package depend

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"go.uber.org/fx"

	l "plassstic.tech/gopkg/golang-manager/internal/depend/logger"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/bot/handler"
)

func NewTelegramBot(handlers []TelegramHandler, lc fx.Lifecycle, config *Config) *bot.Bot {
	log := l.GetLogger("depend.bot")
	var botOptions []bot.Option
	// botOptions = ...
	b, err := bot.New(config.BotToken, botOptions...)
	if err != nil {
		log.Panic(fmt.Sprintf("panic! <%T> %v", err, err))
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go b.Start(context.Background())

				for _, h := range handlers {
					h.Register(b)
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				_, err := b.Close(context.Background())
				return err
			},
		},
	)

	return b
}
