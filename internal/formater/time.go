package formater

import (
	"fmt"
	"strings"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/tj/go-naturaldate"
)

const (
	// DateFormatICal is the date format used by iCal
	DateFormatICal = "20060102T150405Z"
	// DateFormatDefault is the default date format used by remindme
	DateFormatDefault = "15:04 02.01.2006 (MST)"
)

// ToLocalTime converts the time object to a localized time string
func ToLocalTime(datetime time.Time, timezone string) string {
	if timezone == "" {
		return datetime.UTC().Format(DateFormatDefault)
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return datetime.UTC().Format(DateFormatDefault)
	}

	return datetime.In(loc).Format(DateFormatDefault)
}

// ParseTime parses the time into the local timezone of a channel. If no timezone is given it defaults to UTC. Days without a time specified default to 9:00
func ParseTime(msg string, channel *database.Channel, rawDate bool) (time.Time, error) {
	// Clear body from characters the library can not handle
	msg = string(alphaNumericString([]byte(StripReply(msg))))

	loc := time.UTC
	if channel != nil {
		if parsedLoc, err := time.LoadLocation(channel.TimeZone); err == nil {
			loc = parsedLoc
		}
	}
	baseTime := time.Now().In(loc)

	parsedTime, err := naturaldate.Parse(msg, baseTime, naturaldate.WithDirection(naturaldate.Future))
	if err != nil {
		return parsedTime, err
	}

	// Past? then set to in an hour
	if !rawDate && time.Until(parsedTime) <= 5*time.Minute {
		parsedTime = time.Now().Add(time.Hour).In(loc)
	}

	// Midnight? Move to 9:00
	if !rawDate {
		timeString := parsedTime.In(loc).Format("15:04")
		if timeString == "00:00" && !(strings.Contains(msg, "00:00") || strings.Contains(msg, "12am") || strings.Contains(msg, "24:00")) {
			parsedTime = parsedTime.Add(9 * time.Hour)
		}
	}

	return parsedTime.In(loc), nil
}

func alphaNumericString(in []byte) []byte {
	out := make([]byte, len(in))

	for i := range in {
		if (in[i] >= 'a' && in[i] <= 'z') ||
			(in[i] >= 'A' && in[i] <= 'Z') ||
			(in[i] >= '0' && in[i] <= '9') ||
			(in[i] == ':') {
			out[i] = in[i]
		} else {
			out[i] = ' '
		}
	}

	return out
}

// TimeToHourAndMinute converts a time object to an string with the hour and minute in 24h format
func TimeToHourAndMinute(t time.Time) string {
	hours := t.Hour()
	minutes := t.Minute()
	if minutes < 10 {
		return fmt.Sprintf("%d:0%d", hours, minutes)
	}

	return fmt.Sprintf("%d:%d", hours, minutes)
}
