package main

import (
	"github.com/go-telegram/bot"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	. "plassstic.tech/gopkg/golang-manager/internal/depend"
	l "plassstic.tech/gopkg/golang-manager/internal/depend/logger"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/bot"
)

func main() {
	fx.New(
		fx.WithLogger(func() fxevent.Logger { return &fxevent.ZapLogger{Logger: l.GetLogger("zap").Desugar()} }),

		fx.Provide(
			NewConfig,
			NewEntClient,
			fx.Annotate(NewTelegramBot, fx.ParamTags(`group:"handlers"`)),
		),

		Handlers(),

		fx.Invoke(func(*bot.Bot) {}),
	).Run()
}
