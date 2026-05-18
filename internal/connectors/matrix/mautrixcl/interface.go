package mautrixcl

import (
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

// Client wraps mautrix.Client for mocking.
type Client interface {
	JoinedMembers(roomID id.RoomID) (resp *mautrix.RespJoinedMembers, err error)
}
