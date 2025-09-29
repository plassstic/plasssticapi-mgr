package bot

import (
	"context"
	"crypto/sha3"
	"fmt"
	"strings"
	"time"

	tg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/samber/lo"
	"plassstic.tech/gopkg/golang-manager/internal/depend/logger"
	"plassstic.tech/gopkg/golang-manager/internal/logic/api"
	"plassstic.tech/gopkg/golang-manager/internal/logic/bot/utils"
	"plassstic.tech/gopkg/golang-manager/internal/service"
	"plassstic.tech/gopkg/golang-manager/lib/ent"
	"plassstic.tech/gopkg/golang-manager/lib/ent/schema"
)

func True() *bool {
	b := true
	return &b
}

type goroutinesRepo struct {
	D map[int]chan struct{}
}

var Repo *goroutinesRepo

func gr(token string, uid int, e schema.Editable, closeChan <-chan struct{}) {
	log := logger.GetLogger(fmt.Sprintf("thread -> %v", uid))
	bot, err := tg.New(token)
	if err != nil {
		_, _ = bot.SendMessage(context.Background(), &tg.SendMessageParams{ChatID: uid, Text: fmt.Sprintf("failed to start your bot %v", err)})
		return
	}
	var lastHash string
	var kEq int
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-closeChan:
			_, _ = bot.Close(context.Background())
			log.Infof("successfully exited for user %v", uid)
			return
		case <-ticker.C:
			pl, err := api.GetPlayer(int64(uid))

			if err != nil {
				log.Errorf("failed to get player for uid %v (%v)", uid, err)
				continue
			}

			if pl == nil {
				kEq++
				log.Infof("204, sleeping %v secs", 5*kEq)
				time.Sleep(time.Duration(5*kEq) * time.Second)
				continue
			}

			text := utils.GetFormattedPlayer(pl)

			hash := fmt.Sprintf("%x", sha3.Sum256([]byte(text)))

			if hash == lastHash {
				kEq++
				log.Infof("same payload, sleeping %v secs", 5*kEq)
				continue
			} else {
				kEq = 0
				lastHash = hash
			}

			_, err = bot.EditMessageText(context.Background(), &tg.EditMessageTextParams{
				ChatID:    e.ChatId,
				MessageID: e.Id,
				Text:      text,
				ParseMode: models.ParseModeMarkdown,
				LinkPreviewOptions: &models.LinkPreviewOptions{
					IsDisabled: True(),
				},
			})

			if err != nil {
				log.Errorf("failed to edit msg for uid %v (%v)", uid, err)
				continue
			}

			log.Infof("200 OK, %s - %s", pl.Item.Name, strings.Join(lo.Map(pl.Item.Artists, func(item api.ArtistSI, index int) string { return item.Name }), ", "))
		}
	}
}

func Init(c *ent.Client) {
	Repo = &goroutinesRepo{D: make(map[int]chan struct{})}
	Repo.populate(c)
}

func (r *goroutinesRepo) populate(c *ent.Client) {
	eligible, err := (&service.UserService{}).With(c).GetAllNonNil().Many()
	if err != nil {
		panic("failed to get eligible users for start")
	}

	for _, u := range eligible {
		if u.Bot.Token == "" ||
			u.Bot.Handle == "" ||
			u.Editable.Id == 0 ||
			u.Editable.ChatId == 0 {
			continue
		}
		r.D[u.ID] = make(chan struct{}, 1)
		go gr(u.Bot.Token, u.ID, u.Editable, r.D[u.ID])
	}
}

func (r *goroutinesRepo) Stop() {
	for _, closeChan := range r.D {
		closeChan <- struct{}{}
	}
}
