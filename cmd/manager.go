package main

import (
	tg "github.com/go-telegram/bot"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/depend"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/depend/logger"
)

func main() {
	fx.New(
		fx.WithLogger(func() fxevent.Logger { return &fxevent.ZapLogger{Logger: GetLogger("zap").Desugar()} }),

		fx.Provide(
			NewConfig,
			NewEntClient,
			NewTelegramBot,
		),

		fx.Invoke(func(*tg.Bot) {}),
	).Run()
}
