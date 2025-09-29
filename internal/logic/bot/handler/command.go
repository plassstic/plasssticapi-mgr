package handler

import (
	"context"
	"fmt"

	tg "github.com/go-telegram/bot"
	tgm "github.com/go-telegram/bot/models"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/api"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/bot/utils"
	. "plassstic.tech/gopkg/golang-manager/internal/service"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
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
			MFunc:   MSGExact("/start"),
		},
		{
			Handler: c.handleStart,
			MFunc:   CQExact("start_menu"),
		},
	}
}

func (c *cmd) handleStart(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	var text string

	info := GetUInfo(u)

	if info == nil {
		return
	}

	result, err := (&UserService{}).
		With(c.client).
		Ensure(int(info.User.ID)).
		One()

	var kb *tgm.InlineKeyboardMarkup

	if err != nil {
		text = "К сожалению, произошла ошибка при создании аккаунта"
	} else {
		me, err := GetMe(info.User.ID)
		if err != nil {
			text = fmt.Sprintf("К сожалению, произошла ошибка при запросе к API %v", err)
		} else {
			kb = &tgm.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgm.InlineKeyboardButton{
					{
						{Text: "🔧Настройки", CallbackData: "settings"},
					},
				},
			}
			text = "Привет, %s (TG_ID: %d, SP_ID: %s) 👋\n\n" +
				"В этом боте ты сможешь подключить в свой канал автообновляющийся статус Spotify :0"
			text = fmt.Sprintf(text, info.User.Handle, result.ID, me["id"].(string))
		}
	}

	info.Respond(ctx, b, text, kb)
}
