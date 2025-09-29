package handler

import (
	"slices"

	"github.com/go-telegram/bot"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
)

type TelegramHandler struct {
	Handler bot.HandlerFunc
	MFunc   bot.MatchFunc
}

func Register(b *bot.Bot, c *ent.Client) {
	handlers := slices.Concat(cmdHandlers(c), profHandlers(c))
	for _, h := range handlers {
		b.RegisterHandlerMatchFunc(h.MFunc, h.Handler)
	}
}
