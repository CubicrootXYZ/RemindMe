package format

import (
	"fmt"
	"strings"
	"time"

	"github.com/tj/go-naturaldate"

	_ "time/tzdata" // Import timezone data.
)

const (
	// DateTimeFormatDefault is the default date format used by remindme
	DateTimeFormatDefault = "15:04 02.01.2006 (MST)"
	// DateFormatShort is a short date format used.
	DateFormatShort = "Mon, 02 Jan"
	// TimeFormatShot is a short time format used.
	TimeFormatShort = "15:04 (MST)"
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
		return datetime.UTC().Format(DateTimeFormatDefault)
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return datetime.UTC().Format(DateTimeFormatDefault)
	}

	return datetime.In(loc).Format(DateTimeFormatDefault)
}

// ToShortLocalTime converts the time object to a short localized time string
func ToShortLocalTime(datetime time.Time, timezone string) string {
	if timezone == "" {
		return datetime.UTC().Format(TimeFormatShort)
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return datetime.UTC().Format(TimeFormatShort)
	}

	return datetime.In(loc).Format(TimeFormatShort)
}

func toLocalTime(datetime time.Time, loc *time.Location) string {
	return datetime.In(loc).Format(DateTimeFormatDefault)
}

// TimeToHourAndMinute converts a time object to an string with the hour and minute in 24h format.
func TimeToHourAndMinute(t time.Time) string {
	hours := t.Hour()
	minutes := t.Minute()
	if minutes < 10 {
		return fmt.Sprintf("%d:0%d", hours, minutes)
	}

	return fmt.Sprintf("%d:%d", hours, minutes)
}

// ToNiceDuration formats a time.Duration into a nice string.
func ToNiceDuration(d time.Duration) string {
	pre := ""
	if d < 0 {
		d *= -1
		pre = "-"
	}

	if d < time.Minute {
		return fmt.Sprintf("%s%.0f seconds", pre, float64(d/time.Second))
	} else if d < time.Hour {
		return fmt.Sprintf("%s%.0f minutes", pre, float64(d/time.Minute))
	} else if d < 48*time.Hour {
		return fmt.Sprintf("%s%.0f hours", pre, float64(d/time.Hour))
	}
	return fmt.Sprintf("%s%.0f days", pre, float64(d/(24*time.Hour)))
}

func tzFromString(timezone string) *time.Location {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.UTC
	}
	return loc
}
