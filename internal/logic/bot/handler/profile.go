package handler

import (
	"context"
	"fmt"

	tg "github.com/go-telegram/bot"
	tgm "github.com/go-telegram/bot/models"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/bot"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/bot/utils"
	"plassstic.tech/gopkg/golang-manager/internal/service"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
	"plassstic.tech/gopkg/golang-manager/lib/fsm"
)

const (
	fsProfNewToken      fsm.StateID = "newToken"
	fsProfNewMessage    fsm.StateID = "newMessage"
	fsProfChangeMessage fsm.StateID = "changeMessage"
)

type profile struct {
	client *ent.Client
}

func profHandlers(client *ent.Client) []*TelegramHandler {
	p := &profile{
		client: client,
	}

	return []*TelegramHandler{
		{
			Handler: p.handleMenu,
			MFunc:   CQExact("settings"),
		},
		{
			Handler: p.handleToken,
			MFunc:   CQExact("add_bot"),
		},
	}
}

func (p *profile) handleMenu(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	var text string
	var kb *tgm.InlineKeyboardMarkup

	info := GetUInfo(u)
	user, err := (&service.UserService{}).With(p.client).Get(int(info.User.ID)).One()

	if err != nil {
		return
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

	info := GetUInfo(u)

	if info == nil {
		return
	}

	text = "Введите BOT_TOKEN для бота с назначенными правами администратора в целевом канале:\n\n<i>Подсказка: Создать бота можно через @botfather, отправив <b>ему</b> команду <code>/newbot</code></i>"

	_ = GlobalFSM.Transition(ctx, info.User.ID, fsProfNewToken, b, u)

	info.Respond(ctx, b, text, NilMarkup)
}

func (p *profile) handleMessage(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	var text string

	text = "Перешлите сообщение из канала с ботом-администратором, которое вы хотите выделить под Spotify:"

	GetUInfo(u).Respond(ctx, b, text, NilMarkup)
}

func (p *profile) handleClear(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	var text string

	info := GetUInfo(u)

	text = "Перешлите сообщение из канала с ботом-администратором, которое вы хотите выделить под Spotify:"

	info.Respond(ctx, b, text, NilMarkup)
}
