package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/fx"
	"plassstic.tech/gopkg/golang-manager/internal/logic/bot/middleware"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
)

type TelegramHandler interface {
	Register(*bot.Bot)
}
type telegramHandler func(ctx context.Context, bot *bot.Bot, update *models.Update, client *ent.Client)

type handlerWrap struct {
	client      *ent.Client
	handler     telegramHandler
	handlerType bot.HandlerType
	pattern     string
	matchType   bot.MatchType
}

type TelegramHandlerOpts struct {
	Handler bot.HandlerType
	Pattern string
	Match   bot.MatchType
}

func NewTelegramHandler(handler telegramHandler, opts *TelegramHandlerOpts) any {
	return func(client *ent.Client) TelegramHandler {
		return &handlerWrap{
			client: client,

			// handler
			handler: handler,

			// defaults
			handlerType: opts.Handler,
			pattern:     opts.Pattern,
			matchType:   opts.Match,
		}
	}
}

func (c *handlerWrap) applyMiddleware() bot.HandlerFunc {
	wrapped := func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		c.handler(ctx, bot, update, c.client)
	}

	for _, mw := range middleware.Middlewares() {
		wrapped = mw(wrapped)
	}
	return wrapped
}

func (c *handlerWrap) getHandler() bot.HandlerFunc {
	return c.applyMiddleware()
}

func (c *handlerWrap) Register(b *bot.Bot) {
	b.RegisterHandler(
		c.handlerType,
		c.pattern,
		c.matchType,
		c.getHandler(),
	)
}

func AsHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(TelegramHandler)),
		fx.ResultTags(`group:"handlers"`),
	)
}
