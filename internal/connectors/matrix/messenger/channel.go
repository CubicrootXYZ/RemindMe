package asyncmessenger

import (
	"github.com/dchest/uniuri"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

// ChannelResponse holds information about a channel
type ChannelResponse struct {
	ChannelExternalIdentifier string
	UserIdentifier            string
}

// CreateChannel creates a new matrix channel with the given user
func (messenger *service) CreateChannel(userIdentifier string) (*ChannelResponse, error) {
	room := mautrix.ReqCreateRoom{
		Visibility:    "private",
		RoomAliasName: "RemindMe-" + uniuri.NewLen(5),
		Name:          "RemindMe",
		Topic:         "I will be your personal reminder bot",
		Invite:        []id.UserID{id.UserID(userIdentifier)},
		Preset:        "trusted_private_chat",
	}

	response, err := messenger.client.CreateRoom(&room)
	if err != nil {
		return nil, err
	}

	return &ChannelResponse{
		ChannelExternalIdentifier: response.RoomID.String(),
		UserIdentifier:            userIdentifier,
	}, nil
}
