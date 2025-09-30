package utils

import (
	tg "github.com/go-telegram/bot"
	tgm "github.com/go-telegram/bot/models"
)

//func CQPrefix(p string) tg.MatchFunc {
//	return func(u *tgm.Update) bool {
//		return u.CallbackQuery != nil && strings.HasPrefix(u.CallbackQuery.Data, p)
//	}
//}

func CQExact(p string) tg.MatchFunc {
	return func(u *tgm.Update) bool {
		return u.CallbackQuery != nil && u.CallbackQuery.Data == p
	}
}
func MSGExact(p string) tg.MatchFunc {
	return func(u *tgm.Update) bool {
		return u.Message != nil && u.Message.Text == p
	}
}
