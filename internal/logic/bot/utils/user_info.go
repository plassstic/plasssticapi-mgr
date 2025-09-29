package utils

import "github.com/go-telegram/bot/models"

type userInfo struct {
	ID     int64
	Handle string
	ChatID int64
}

type UpdateInfo struct {
	User    userInfo
	Payload string
	Type    interface{}
	Msg     *models.Message
}

func GetUInfo(u *models.Update) *UpdateInfo {
	switch {
	case u.Message != nil:
		return &UpdateInfo{
			User: userInfo{
				ID:     u.Message.From.ID,
				Handle: u.Message.From.Username,
				ChatID: u.Message.Chat.ID,
			},
			Payload: u.Message.Text,
			Type:    models.Message{},
			Msg:     u.Message,
		}
	case u.CallbackQuery != nil:
		return &UpdateInfo{
			User: userInfo{
				ID:     u.CallbackQuery.From.ID,
				Handle: u.CallbackQuery.From.Username,
				ChatID: u.CallbackQuery.Message.Message.Chat.ID,
			},
			Payload: u.CallbackQuery.Data,
			Type:    models.CallbackQuery{},
			Msg:     u.CallbackQuery.Message.Message,
		}
	default:
		return nil
	}
}
