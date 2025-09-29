package handler

import (
	"context"

	tg "github.com/go-telegram/bot"
	tgm "github.com/go-telegram/bot/models"
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
		text += "\n\n" + "К вашему аккаунту подключен бот @%s"
		kb = &tgm.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgm.InlineKeyboardButton{
				{
					{Text: "🤖 Сменить бота", CallbackData: "add_bot"},
					{Text: "💬 Изменить сообщение", CallbackData: "change_message"},
				},
				{
					{Text: "🔙 Назад", CallbackData: "start_menu"},
				},
			},
		}
	}

	info.Respond(ctx, b, text, kb)
}

func (p *profile) handleToken(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	var text string

	info := GetUInfo(u)

	if info == nil {
		return
	}

	text = "Введите токен бота:"

	_ = GlobalFSM.Transition(ctx, info.User.ID, fsProfNewToken, b, u)

	info.Respond(ctx, b, text, &tgm.InlineKeyboardMarkup{InlineKeyboard: [][]tgm.InlineKeyboardButton{}})
}

func (p *profile) handleMessage(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	var text string

	text = "Перешлите сообщение из канала, где бот назначен администратором с правом управлять сообщениями, которое вы хотите выделить под Spotify:"

	GetUInfo(u).Respond(ctx, b, text, &tgm.InlineKeyboardMarkup{InlineKeyboard: [][]tgm.InlineKeyboardButton{}})
}
