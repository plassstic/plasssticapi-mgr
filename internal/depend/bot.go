package depend

import (
	"context"
	"fmt"

	tg "github.com/go-telegram/bot"
	"go.uber.org/fx"
	l "plassstic.tech/gopkg/golang-manager/internal/depend/logger"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/bot/handler"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/bot/middleware"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
)

func NewTelegramBot(lc fx.Lifecycle, config *Config, c *ent.Client) *tg.Bot {
	log := l.GetLogger("depend.bot")

	botOptions := []tg.Option{tg.WithMiddlewares(Middlewares()...), tg.WithDebug(), tg.WithDefaultHandler(FSM(c))}
	b, err := tg.New(config.BotToken, botOptions...)
	if err != nil {
		log.Panic(fmt.Sprintf("panic! <%Type> %v", err, err))
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				Register(b, c)

				go b.Start(context.Background())

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
