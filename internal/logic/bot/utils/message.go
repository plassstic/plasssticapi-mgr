package utils

import (
	"context"

	tg "github.com/go-telegram/bot"
	tgm "github.com/go-telegram/bot/models"
	"plassstic.tech/gopkg/golang-manager/internal/depend/logger"
)

func sendMessage(ctx context.Context, bot *tg.Bot, chatID int64, text string, replyMarkup tgm.ReplyMarkup) {
	log := logger.GetLogger("logic -> bot -> utils")

	params := &tg.SendMessageParams{ChatID: chatID, Text: text, ParseMode: tgm.ParseModeHTML, ReplyMarkup: replyMarkup}

	msg, err := bot.SendMessage(ctx, params)

	if err != nil {
		log.Errorf("error! failure to send msg to user %d; %v", chatID, err)
	} else {
		log.Infof("sent msg id %d to user cid %d", msg.ID, chatID)
	}
}

func safeEditMessage(ctx context.Context, bot *tg.Bot, chatID int64, messageID int, text string, replyMarkup tgm.ReplyMarkup) {
	log := logger.GetLogger("logic -> bot -> utils")

	msg, err := bot.EditMessageText(ctx, &tg.EditMessageTextParams{ChatID: chatID, MessageID: messageID, Text: text, ParseMode: tgm.ParseModeHTML, ReplyMarkup: replyMarkup})

	if err != nil {
		sendMessage(ctx, bot, chatID, text, replyMarkup)
	} else {
		log.Infof("edited msg id %d for user %d", msg.ID, msg.Chat.ID)
	}
}

func (u *UpdateInfo) Respond(ctx context.Context, bot *tg.Bot, text string, replyMarkup tgm.ReplyMarkup) {
	switch u.Type.(type) {
	case tgm.Message:
		sendMessage(ctx, bot, u.User.ChatID, text, replyMarkup)
	case tgm.CallbackQuery:
		safeEditMessage(ctx, bot, u.User.ChatID, u.Msg.ID, text, replyMarkup)
	}
}
