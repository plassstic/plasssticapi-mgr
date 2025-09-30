package handler

import (
	"context"
	"fmt"

	tg "github.com/go-telegram/bot"
	tgm "github.com/go-telegram/bot/models"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/api"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/bot"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/bot/utils"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/thread"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/service"
	"plassstic.tech/gopkg/plassstic-mgr/lib/ent"
)

type cmd struct {
	client *ent.Client
}

func cmdHandlers(client *ent.Client) []*TelegramHandler {
	c := &cmd{
		client: client,
	}
	return []*TelegramHandler{
		{
			Handler: c.handleStart,
			Match:   MSGExact("/start"),
		},
		{
			Handler: c.handleStart,
			Match:   CQExact("start_menu"),
		},
	}
}

func (c *cmd) handleStart(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	var text string

	info := UIFromUpdate(u)

	if info == nil {
		return
	}

	user, err := (&UserService{}).
		With(c.client).
		Ensure(int(info.User.ID)).
		One()

	kb := NilMarkup

	if err != nil {
		text = fmt.Sprintf("Произошла ошибка при создании аккаунта: %v", err)
	} else {
		me, err := GetMe(info.User.ID)
		if err != nil {
			text = fmt.Sprintf("Произошла ошибка при запросе к API: %v\n\nЕсли вы еще не связывали свой аккаунт Telegram с Spotify, это можно сделать здесь: https://api.plassstic.tech/public/auth", err)
		} else {
			kb = &tgm.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgm.InlineKeyboardButton{
					{
						{Text: "🔧Настройки", CallbackData: "settings"},
					},
				},
			}
			text = "Привет, %s (TG_ID: %d, SP_ID: %s) 👋\n\n" +
				"В этом боте ты сможешь подключить в свой канал автообновляющийся статус Spotify :)"

			if HasThread(info.User.ID) {
				text += fmt.Sprintf("\n\n✅ Активен поток для бота @%s", user.Bot.Handle)
			}
			text = fmt.Sprintf(text, info.User.Handle, user.ID, me["id"])
		}
	}

	info.Respond(ctx, b, text, kb)
}
