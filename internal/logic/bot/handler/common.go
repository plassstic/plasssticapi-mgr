package handler

import (
	"slices"

	tg "github.com/go-telegram/bot"
	"plassstic.tech/gopkg/plassstic-mgr/lib/ent"
)

type TelegramHandler struct {
	Handler tg.HandlerFunc
	Match   tg.MatchFunc
}

func Register(b *tg.Bot, c *ent.Client) {
	handlers := slices.Concat(cmdHandlers(c), profileHandlers(c))
	for _, h := range handlers {
		b.RegisterHandlerMatchFunc(h.Match, h.Handler)
	}
}
