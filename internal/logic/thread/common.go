package bot

import (
	"context"
	"crypto/sha3"
	"fmt"
	"strings"
	"time"

	tg "github.com/go-telegram/bot"
	tgm "github.com/go-telegram/bot/models"
	"github.com/samber/lo"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/depend/logger"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/api"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/bot/utils"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/service"
	"plassstic.tech/gopkg/plassstic-mgr/lib/ent"
	"plassstic.tech/gopkg/plassstic-mgr/lib/ent/schema"
)

func True() *bool {
	b := true
	return &b
}

type goroutinesRepo struct {
	threads map[int]chan struct{}
}

var repo *goroutinesRepo

func runThread(userId int, token string, editable schema.Editable, closeChan <-chan struct{}) {
	var bot *tg.Bot
	var err error

	log := GetLogger(fmt.Sprintf("thread.%v", userId))

	if bot, err = tg.New(token); err != nil {
		_, _ = bot.SendMessage(context.Background(), &tg.SendMessageParams{ChatID: userId, Text: fmt.Sprintf("failed to start your bot %v", err)})
		return
	}

	var lastHash string
	var backoffFactor int

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-closeChan:
			bot.Close(context.Background())
			log.Infof("successfully exited for user %v", userId)

			return
		case <-ticker.C:
			pl, err := GetPlayer(int64(userId))

			if err != nil {
				log.Errorf("failed to get player for userId %v (%e)", userId, err)
				continue
			}

			if pl == nil {
				backoffFactor++

				delay := min(5*backoffFactor, 120)
				log.Infof("204 No Content, player is inactive, delay %vs", delay)
				time.Sleep(time.Duration(delay) * time.Second)

				continue
			} else {
				backoffFactor = 0
			}

			text := GetFormattedPlayer(pl)
			hash := fmt.Sprintf("%x", sha3.Sum256([]byte(text)))

			if hash == lastHash {
				backoffFactor++

				delay := min(5*backoffFactor, 120)
				log.Infof("200 OK, same payload, delay %vs", delay)
				time.Sleep(time.Duration(delay) * time.Second)

				continue
			} else {
				backoffFactor = 0
				lastHash = hash
			}

			_, err = bot.EditMessageText(context.Background(), &tg.EditMessageTextParams{
				ChatID:    editable.ChatId,
				MessageID: editable.Id,
				Text:      text,
				ParseMode: tgm.ParseModeMarkdown,
				LinkPreviewOptions: &tgm.LinkPreviewOptions{
					IsDisabled: True(),
				},
			})

			if err != nil {
				log.Errorf("failed to edit msg for userId %v (%e)", userId, err)
				continue
			}

			log.Infof("200 OK, set presence to %s - %s", pl.Item.Name, strings.Join(lo.Map(pl.Item.Artists, func(item ArtistSI, index int) string { return item.Name }), ", "))
		}
	}
}

func (repo *goroutinesRepo) populate(c *ent.Client) {
	eligible, err := (&UserService{}).
		With(c).
		GetAll().
		FilterMany(
			func(user *ent.User, _ int) bool {
				p := user.Bot.Token == "" ||
					user.Bot.Handle == "" ||
					user.Editable.Id == 0 ||
					user.Editable.ChatId == 0
				return !p
			},
		).
		Many()
	if err != nil {
		GetLogger("thread.populate").Panicf("panic! failed to get eligible users for start (%e)", err)
	}

	for _, user := range eligible {
		repo.threads[user.ID] = make(chan struct{}, 1)
		go runThread(user.ID, user.Bot.Token, user.Editable, repo.threads[user.ID])
	}
}
