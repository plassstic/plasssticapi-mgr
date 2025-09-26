package main

import (
	"github.com/go-telegram/bot"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	. "plassstic.tech/gopkg/golang-manager/internal/depend"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/bot"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.SugaredLogger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Desugar().Named("zap")}
		}),

		fx.Provide(
			NewLogger,
			NewConfig,
			NewEntClient,
			fx.Annotate(NewTelegramBot, fx.ParamTags(`group:"handlers"`)),
		),

		Handlers(),

		fx.Invoke(func(*bot.Bot) {}),
	).Run()
}
