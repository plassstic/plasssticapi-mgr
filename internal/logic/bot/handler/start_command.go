package handler

import (
	"context"
	"fmt"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	. "plassstic.tech/gopkg/golang-manager/internal/depend/logger"
	. "plassstic.tech/gopkg/golang-manager/internal/service"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
)

func StartCommand(ctx context.Context, bot *tg.Bot, update *models.Update, client *ent.Client) {
	log := GetLogger("bot.handlers.StartCommand")

	log.Info("catch")

	var result interface{}

	result = (&UserService{}).
		With(client).
		Ensure(int(update.Message.From.ID)).
		Fin()

	switch result.(type) {
	case error:
		_, err := bot.SendMessage(ctx, &tg.SendMessageParams{ChatID: update.Message.From.ID, Text: "Failed"})
		if err != nil {
			log.Errorf("error! %T: %v", err, err)
		}
	case *ent.User:
		_, err := bot.SendMessage(ctx, &tg.SendMessageParams{ChatID: update.Message.From.ID, Text: fmt.Sprintf("OK %#v", result.(*ent.User))})
		if err != nil {
			log.Errorf("error! %T: %v", err, err)
		}
	default:
		log.Infof("%#v", result)
	}

}
