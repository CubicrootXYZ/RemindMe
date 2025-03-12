package matrix

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/daemon"
)

func (service *service) ToLocalTime(date time.Time, output *daemon.Output) time.Time {
	room, err := service.matrixDatabase.GetRoomByID(output.OutputID)
	if err != nil {
		return date
	}

	if room.TimeZone == "" {
		return date
	}

	loc, err := time.LoadLocation(room.TimeZone)
	if err != nil {
		return date
	}

	return date.In(loc)
}
