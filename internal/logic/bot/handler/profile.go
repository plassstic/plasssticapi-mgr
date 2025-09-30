package handler

import (
	"context"
	"fmt"

	tg "github.com/go-telegram/bot"
	tgm "github.com/go-telegram/bot/models"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/depend/logger"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/bot"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/bot/handler/state"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/bot/utils"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/service"
	"plassstic.tech/gopkg/plassstic-mgr/lib/ent"
	"plassstic.tech/gopkg/plassstic-mgr/lib/fsm"
)

type profile struct {
	client *ent.Client
}

func profileHandlers(client *ent.Client) []*TelegramHandler {
	p := &profile{
		client: client,
	}

	p.setProfileCallbacks()

	return []*TelegramHandler{
		{
			Handler: p.handleMenu,
			Match:   CQExact("settings"),
		},
		{
			Handler: p.handleToken,
			Match:   CQExact("add_bot"),
		},
	}
}

func (p *profile) setProfileCallbacks() {
	GlobalFSM.AddCallbacks(map[fsm.StateID]fsm.Callback{
		ProfileNewTokenState: func(ctx context.Context, b *tg.Bot, update *tgm.Update) {
			return
		},
		ProfileNewMessageState: p.handleMessage,
		ProfileChangeMessageState: func(ctx context.Context, b *tg.Bot, update *tgm.Update) {
			return
		},
	})
}

func (p *profile) handleMenu(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	var text string
	var kb *tgm.InlineKeyboardMarkup

	info := UIFromUpdate(u)
	user, err := (&UserService{}).
		With(p.client).
		Get(int(info.User.ID)).
		One()

	if err != nil {
		GetLogger("profile.handleMenu").Panicf("panic! %e", err)
	}

	text = "🛠 Настройки"

	if user.Bot.Token == "" {
		kb = &tgm.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgm.InlineKeyboardButton{
				{
					{Text: "🤖 Добавить бота", CallbackData: "add_bot"},
				},
				{
					{Text: "🔙 Назад", CallbackData: "start_menu"},
				},
			},
		}
	} else {
		text += "\n\n" + fmt.Sprintf("К вашему аккаунту подключен бот @%s", user.Bot.Handle)
		kb = &tgm.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgm.InlineKeyboardButton{
				{
					{Text: "🤖 Очистить", CallbackData: "clear"},
				},
				{
					{Text: "🔙 Назад", CallbackData: "start_menu"},
				},
			},
		}
	} // TODO(plassstic): Вынести клавиатуры в отд утилиту

	info.Respond(ctx, b, text, kb)
}

func (p *profile) handleToken(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	var text string

	info := UIFromUpdate(u)

	text = "Введите BOT_TOKEN для бота с назначенными правами администратора в целевом канале:\n\n<i>Подсказка: Создать бота можно через @botfather, отправив <b>ему</b> команду <code>/newbot</code></i>"

	_ = GlobalFSM.Transition(ctx, info.User.ID, ProfileNewTokenState, b, u)

	info.Respond(ctx, b, text, NilMarkup)
}

func (p *profile) handleMessage(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	var text string

	text = "Перешлите сообщение из канала с ботом-администратором, которое вы хотите выделить под Spotify:"

	UIFromUpdate(u).Respond(ctx, b, text, NilMarkup)
}

func (p *profile) handleClear(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	var text string

	info := UIFromUpdate(u)

	text = "Перешлите сообщение из канала с ботом-администратором, которое вы хотите выделить под Spotify:"

	info.Respond(ctx, b, text, NilMarkup)
}
