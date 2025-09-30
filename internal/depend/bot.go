package depend

import (
	"context"

	tg "github.com/go-telegram/bot"
	"go.uber.org/fx"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/depend/logger"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/bot/handler"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/bot/handler/state"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/bot/middleware"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/thread"
	"plassstic.tech/gopkg/plassstic-mgr/lib/ent"
)

func NewTelegramBot(lc fx.Lifecycle, config *Config, c *ent.Client) *tg.Bot {
	var bot *tg.Bot
	var err error

	log := GetLogger("depend.spawnBot")
	botOptions := []tg.Option{tg.WithMiddlewares(Middlewares()...), tg.WithDefaultHandler(FSM(c))}

	if bot, err = tg.New(config.BotToken, botOptions...); err != nil {
		log.Panicf("panic! <%T> %v", err, err)
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				Register(bot, c)
				InitPresenceThreads(c)
				go bot.Start(context.Background())

				return nil
			},
			OnStop: func(ctx context.Context) error {
				StopPresenceThreads()
				_, err = bot.Close(context.Background())

				return err
			},
		},
	)

	return bot
}
