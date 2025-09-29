package handler

import (
	"context"
	"regexp"

	tg "github.com/go-telegram/bot"
	tgm "github.com/go-telegram/bot/models"
	"plassstic.tech/gopkg/golang-manager/internal/depend/logger"
	. "plassstic.tech/gopkg/golang-manager/internal/logic/bot/utils"
	. "plassstic.tech/gopkg/golang-manager/internal/service"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
	"plassstic.tech/gopkg/golang-manager/lib/ent/schema"
	"plassstic.tech/gopkg/golang-manager/lib/fsm"
)

const (
	fsUnknown fsm.StateID = "unknown"
)

var GlobalFSM *fsm.FSM[string, string]

type fsmRouter struct {
	client *ent.Client
}

func FSM(c *ent.Client) tg.HandlerFunc {
	fr := &fsmRouter{
		client: c,
	}

	p := &profile{client: c}

	GlobalFSM = fsm.New[string, string](
		fsUnknown,
		map[fsm.StateID]fsm.Callback{
			fsProfNewToken: func(ctx context.Context, b *tg.Bot, update *tgm.Update) {
				return
			},
			fsProfNewMessage: p.handleMessage,
			fsProfChangeMessage: func(ctx context.Context, b *tg.Bot, update *tgm.Update) {
				return
			},
		},
	)

	return fr.handle
}

func (r *fsmRouter) handle(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	info := GetUInfo(u)
	if info == nil {
		return
	}
	log := logger.GetLogger("FSM -> routing")

	switch st, _ := GlobalFSM.Current(info.User.ID); st {
	case fsUnknown:
		return
	case fsProfNewToken:
		r.validateProfNewToken(ctx, b, u, info)
	case fsProfNewMessage:
		r.validateProfNewMessage(ctx, b, u, info)
	case fsProfChangeMessage:
		return
	default:
		log.Errorf("unexpected state %s\n", st)
	}

	return

}

func (r *fsmRouter) validateProfNewToken(ctx context.Context, b *tg.Bot, u *tgm.Update, info *UpdateInfo) {
	token := info.Payload
	re, err := regexp.Compile(`[0-9]+:[A-Za-z_\-0-9]{35}`)

	if err != nil || !re.MatchString(token) {
		info.Respond(ctx, b, "Неверный токен; попробуйте еще раз.", &tgm.InlineKeyboardMarkup{InlineKeyboard: [][]tgm.InlineKeyboardButton{}})
		return
	}

	_, err = tg.New(token)

	if err != nil {
		info.Respond(ctx, b, "Неверный токен; попробуйте еще раз.", &tgm.InlineKeyboardMarkup{InlineKeyboard: [][]tgm.InlineKeyboardButton{}})
		return
	}

	_ = GlobalFSM.Set(info.User.ID, "token", token)
	_ = GlobalFSM.Transition(ctx, info.User.ID, fsProfNewMessage, b, u)
	return
}

func (r *fsmRouter) validateProfNewMessage(ctx context.Context, b *tg.Bot, u *tgm.Update, info *UpdateInfo) {
	fo := u.Message.ForwardOrigin
	if fo == nil || fo.Type != tgm.MessageOriginTypeChannel {
		info.Respond(ctx, b, "Сообщение не является пересланным из канала; попробуйте еще раз.", &tgm.InlineKeyboardMarkup{InlineKeyboard: [][]tgm.InlineKeyboardButton{}})
		return
	}

	token, err := GlobalFSM.Get(info.User.ID, "token")

	if err != nil {
		info.Respond(ctx, b, "Произошла ошибка, токен не найден в вашем состоянии\n\nВведите токен еще раз:", &tgm.InlineKeyboardMarkup{InlineKeyboard: [][]tgm.InlineKeyboardButton{}})
		_ = GlobalFSM.Transition(ctx, info.User.ID, fsProfNewToken, b, u)
		return
	}

	tb, err := tg.New(token)

	if err != nil {
		info.Respond(ctx, b, "Неверный токен\n\nВведите токен еще раз:", &tgm.InlineKeyboardMarkup{InlineKeyboard: [][]tgm.InlineKeyboardButton{}})
		_ = GlobalFSM.Transition(ctx, info.User.ID, fsProfNewToken, b, u)
		return
	}

	_, err = tb.EditMessageText(ctx, &tg.EditMessageTextParams{
		ChatID:    fo.MessageOriginChannel.Chat.ID,
		MessageID: fo.MessageOriginChannel.MessageID,
		Text:      "TestingHash" + tg.RandomString(8),
	})

	if err != nil {
		info.Respond(ctx, b, "Не удалось воспроизвести тестовое изменение сообщения.\n\n"+
			"Убедитесь, что менеджеру выданы все права администратора и перешлите сообщение еще раз",
			&tgm.InlineKeyboardMarkup{InlineKeyboard: [][]tgm.InlineKeyboardButton{}})
		return
	}

	editable := schema.Editable{
		Id:     int(fo.MessageOriginChannel.Chat.ID),
		ChatId: fo.MessageOriginChannel.MessageID,
	}

	tbInfo, _ := tb.GetMe(ctx)
	tbModel := schema.Bot{
		Token:  token,
		Handle: tbInfo.Username,
	}

	_, err = (&UserService{}).
		With(r.client).
		SetBoth(info.User.ID, tbModel, editable).
		One()

	if err != nil {
		info.Respond(ctx, b, "Не удалось сохранить изменения.\n\nПопробуйте еще раз", &tgm.InlineKeyboardMarkup{InlineKeyboard: [][]tgm.InlineKeyboardButton{}})
		_ = GlobalFSM.Transition(ctx, info.User.ID, fsProfNewToken, b, u)
		return
	}

	info.Respond(ctx, b, "OK", &tgm.InlineKeyboardMarkup{InlineKeyboard: [][]tgm.InlineKeyboardButton{}})
	_ = GlobalFSM.Reset(info.User.ID)
	return
}
