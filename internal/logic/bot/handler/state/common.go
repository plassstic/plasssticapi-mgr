package state

import (
	"context"

	tg "github.com/go-telegram/bot"
	tgm "github.com/go-telegram/bot/models"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/depend/logger"
	. "plassstic.tech/gopkg/plassstic-mgr/internal/logic/bot/utils"
	"plassstic.tech/gopkg/plassstic-mgr/lib/ent"
	"plassstic.tech/gopkg/plassstic-mgr/lib/fsm"
)

const (
	unknownState fsm.StateID = "unknownState"
)

var GlobalFSM *fsm.FSM[string, string]

type fsmRouter struct {
	client *ent.Client
}

func FSM(c *ent.Client) tg.HandlerFunc {
	router := &fsmRouter{
		client: c,
	}

	GlobalFSM = fsm.New[string, string](
		unknownState,
		map[fsm.StateID]fsm.Callback{},
	)

	return router.handle
}

func (r *fsmRouter) handle(ctx context.Context, b *tg.Bot, u *tgm.Update) {
	info := UIFromUpdate(u)

	log := GetLogger("fsm.routing")

	switch st, _ := GlobalFSM.Current(info.User.ID); st {
	case unknownState:
		return
	case ProfileNewTokenState:
		r.callbackProfileNewToken(ctx, b, u, info)
	case ProfileNewMessageState:
		r.callbackProfileNewMessage(ctx, b, u, info)
	case ProfileChangeMessageState:
		return
	default:
		log.Errorf("unexpected state %s\n", st)
	}

	return

}
