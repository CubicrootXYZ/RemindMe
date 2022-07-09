package matrixsyncer

import (
	"errors"
	"strconv"
	"strings"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/roles"
)

// createChannel creates a new matrix channel
func (s *Syncer) createChannel(userID string, role roles.Role) (*database.Channel, error) {
	roomCreated, err := s.messenger.CreateChannel(userID)
	if err != nil {
		return nil, err
	}

	return s.daemon.Database.AddChannel(userID, roomCreated.RoomID.String(), role)
}

// getSuffixInt returns a suffixed integer in the given string value
func getSuffixInt(value string) (int, error) {
	splitUp := strings.Split(value, "")
	if len(splitUp) == 0 {
		return 0, errors.New("empty string does not contain integer")
	}

	integerString := splitUp[len(splitUp)-1]

	return strconv.Atoi(integerString)
}
