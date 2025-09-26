package depend

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"go.uber.org/fx"
	"go.uber.org/zap"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/bot/handler"
)

func NewTelegramBot(handlers []TelegramHandler, lc fx.Lifecycle, log *zap.SugaredLogger, config *Config) *bot.Bot {
	var botOptions []bot.Option
	// botOptions = ...
	b, err := bot.New(config.BotToken, botOptions...)
	if err != nil {
		log.Named("depend.bot").Panic(fmt.Sprintf("panic! <%T> %v", err, err))
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go b.Start(ctx)

				for _, h := range handlers {
					h.Register(b)
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				_, err := b.Close(ctx)
				return err
			},
		},
	)

	return b
}
