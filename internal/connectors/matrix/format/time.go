package format

import (
	"strings"
	"time"

	"github.com/tj/go-naturaldate"
)

const (
	// DateFormatDefault is the default date format used by remindme
	DateFormatDefault = "15:04 02.01.2006 (MST)"
)

// ParseTime parses the time from the input.
// If a timezone is given the returned time.Time will be in that timezone.
// rawDate disabled will try to find a date in the future.
func ParseTime(msg string, timeZone string, rawDate bool) (time.Time, error) {
	// Clear body from characters the library can not handle
	msg = string(alphaNumericString([]byte(StripReply(msg))))

	loc := time.UTC
	if parsedLoc, err := time.LoadLocation(timeZone); err == nil {
		loc = parsedLoc
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
