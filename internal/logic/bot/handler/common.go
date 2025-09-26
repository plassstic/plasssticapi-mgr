package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
)

type TelegramHandler interface {
	Register(*bot.Bot)
}

type telegramHandler func(ctx context.Context, bot *bot.Bot, update *models.Update, client *ent.Client, log *zap.SugaredLogger)

type handlerWrap struct {
	client      *ent.Client
	log         *zap.SugaredLogger
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

func NewTelegramHandler(handler telegramHandler, opts *TelegramHandlerOpts) func(client *ent.Client, log *zap.SugaredLogger) TelegramHandler {
	return func(client *ent.Client, log *zap.SugaredLogger) TelegramHandler {
		return &handlerWrap{
			client: client,
			log:    log,

			// handler
			handler: handler,

			// defaults
			handlerType: opts.Handler,
			pattern:     opts.Pattern,
			matchType:   opts.Match,
		}
	}
}

func (c *handlerWrap) Register(b *bot.Bot) {
	b.RegisterHandler(
		c.handlerType,
		c.pattern,
		c.matchType,
		func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			c.handler(ctx, bot, update, c.client, c.log)
		},
	)
}

func AsHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(TelegramHandler)),
		fx.ResultTags(`group:"handlers"`),
	)
}
