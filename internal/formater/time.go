package formater

import (
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

// ToLocalTime converts the time object to a localized time string
func ToLocalTime(datetime time.Time, channel *database.Channel) string {
	if channel == nil || channel.TimeZone == "" {
		return datetime.UTC().Format("15:04 02.01.2006 (MST)")
	}

	loc, err := time.LoadLocation(channel.TimeZone)
	if err != nil {
		return datetime.UTC().Format("15:04 02.01.2006 (MST)")
	}

	return datetime.In(loc).Format("15:04 02.01.2006 (MST)")
}
