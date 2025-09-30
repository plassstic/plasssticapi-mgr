package state

import (
	"context"
	"regexp"

	tg "github.com/go-telegram/bot"
	tgm "github.com/go-telegram/bot/models"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/bot"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/bot/utils"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/service"
	"plassstic.tech/gopkg/plassstic-mgr/lib/ent/schema"
)

func (r *fsmRouter) callbackProfileNewToken(ctx context.Context, b *tg.Bot, u *tgm.Update, info *UpdateInfo) {
	token := info.Payload
	re, err := regexp.Compile(`[0-9]+:[A-Za-z_\-0-9]{35}`)

	if err != nil || !re.MatchString(token) {
		info.Respond(ctx, b, "Неверный токен; попробуйте еще раз.", NilMarkup)
		return
	}

	_, err = tg.New(token)

	if err != nil {
		info.Respond(ctx, b, "Неверный токен; попробуйте еще раз.", NilMarkup)
		return
	}

	_ = GlobalFSM.Set(info.User.ID, "token", token)
	_ = GlobalFSM.Transition(ctx, info.User.ID, ProfileNewMessageState, b, u)
	return
}

func (r *fsmRouter) callbackProfileNewMessage(ctx context.Context, b *tg.Bot, u *tgm.Update, info *UpdateInfo) {
	fo := u.Message.ForwardOrigin
	if fo == nil || fo.Type != tgm.MessageOriginTypeChannel {
		info.Respond(ctx, b, "Сообщение не является пересланным из канала; попробуйте еще раз.", NilMarkup)
		return
	}

	token, err := GlobalFSM.Get(info.User.ID, "token")

	if err != nil {
		info.Respond(ctx, b, "Произошла ошибка, токен не найден в вашем состоянии\n\nВведите токен еще раз:", NilMarkup)
		_ = GlobalFSM.Transition(ctx, info.User.ID, ProfileNewTokenState, b, u)
		return
	}

	tb, err := tg.New(token)

	if err != nil {
		info.Respond(ctx, b, "Неверный токен\n\nВведите токен еще раз:", NilMarkup)
		_ = GlobalFSM.Transition(ctx, info.User.ID, ProfileNewTokenState, b, u)
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
			NilMarkup)
		return
	}

	editable := schema.Editable{
		Id:     fo.MessageOriginChannel.MessageID,
		ChatId: int(fo.MessageOriginChannel.Chat.ID),
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
		info.Respond(ctx, b, "Не удалось сохранить изменения.\n\nПопробуйте еще раз", NilMarkup)
		_ = GlobalFSM.Transition(ctx, info.User.ID, ProfileNewTokenState, b, u)
		return
	}

	info.Respond(ctx, b, "OK", NilMarkup)
	_ = GlobalFSM.Reset(info.User.ID)
	return
}
