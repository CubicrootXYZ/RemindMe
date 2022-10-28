package matrixsyncer

import (
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
)

// createChannel creates a new matrix channel
func (s *Syncer) createChannel(userID string, role roles.Role) (*database.Channel, error) {
	roomCreated, err := s.messenger.CreateChannel(userID)
	if err != nil {
		return nil, err
	}

	return s.daemon.Database.AddChannel(userID, roomCreated.ChannelExternalIdentifier, role)
}
