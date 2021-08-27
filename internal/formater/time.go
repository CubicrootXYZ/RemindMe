package formater

import (
	"strings"
	"time"
	"unicode"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/tj/go-naturaldate"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

// ParseTime parses the time into the local timezone of a channel. If no timezone is given it defaults to UTC. Days without a time specified default to 9:00
func ParseTime(msg string, channel *database.Channel) (time.Time, error) {
	// Clear body from characters the library can not handle
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	strippedBody, _, err := transform.String(t, StripReply(msg))
	if err == nil {
		msg = strippedBody
	}

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
	if time.Until(parsedTime) <= 5*time.Minute {
		parsedTime = time.Now().Add(time.Hour).In(loc)
	}

	// Midnight? Move to 9:00
	timeString := parsedTime.In(loc).Format("15:04")
	if timeString == "00:00" && !(strings.Contains(msg, "00:00") || strings.Contains(msg, "12am") || strings.Contains(msg, "24:00")) {
		parsedTime = parsedTime.Add(9 * time.Hour)
	}

	return parsedTime, nil
}
