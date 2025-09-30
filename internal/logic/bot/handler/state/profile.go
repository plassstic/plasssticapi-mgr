package state

import "plassstic.tech/gopkg/plassstic-mgr/lib/fsm"

const (
	ProfileNewTokenState      fsm.StateID = "newToken"
	ProfileNewMessageState    fsm.StateID = "newMessage"
	ProfileChangeMessageState fsm.StateID = "changeMessage"
)
