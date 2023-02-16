package format_test

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTime(t *testing.T) {
	testCases := make(map[string]string)
	testCases["UTC"] = "tomorrow 11:45"
	testCases[""] = "tomorrow 11:45"
	testCases["abcdefg"] = "tomorrow at 11:45"
	testCases["Asia/Jakarta"] = "tomorrow 18:45"
	testCases["Asia/Jakarta"] = "tomorrow 18:45?"
	testCases["Asia/Jakarta"] = "tomorrow 18:45!"
	testCases["Asia/Jakarta"] = "tomorrow 18:45="
	testCases["Asia/Jakarta"] = "tomorrow 18:45"
	testCases["Asia/Jakarta"] = "reminder for? tomorrow 18:45"
	testCases["Asia/Jakarta"] = "reminder for 29? tomorrow 18:45"
	testCases["Asia/Jakarta"] = "reminder for tomorrow 18:45"
	testCases["Asia/Jakarta"] = "reminder for tomorrow 01.01.2022 18:45"
	testCases["Asia/Jakarta"] = "reminder for tomorrow 01-02-2022 18:45"

	for timeZone, msg := range testCases {
		is, err := format.ParseTime(msg, timeZone, false)

		require.NoError(t, err, "Can not parse "+msg+" / "+timeZone)
		assert.Equal(t, "11:45", is.UTC().Format("15:04"), "Wrong date from "+msg+" / "+timeZone)
	}
}

func TestParseTimeWithFailure(t *testing.T) {
	testCases := []string{
		"99:00:",
		"99:00:56678:",
		":.",
		"tomorrow: 19:45",
	}

	for _, msg := range testCases {
		_, err := format.ParseTime(msg, "", false)

		assert.Error(t, err, "Should not parse "+msg)
	}
}

func TestToLocalTime(t *testing.T) {
	refTime := refTime()

	testCases := make(map[string]string)
	testCases["UTC"] = "11:45 12.11.2014 (UTC)"
	testCases[""] = "11:45 12.11.2014 (UTC)"
	testCases["abcdefg"] = "11:45 12.11.2014 (UTC)"
	testCases["Europe/Berlin"] = "12:45 12.11.2014 (CET)"
	testCases["America/Mexico_City"] = "05:45 12.11.2014 (CST)"
	testCases["Asia/Jakarta"] = "18:45 12.11.2014 (WIB)"

	for timeZone, should := range testCases {
		is := format.ToLocalTime(refTime, timeZone)

		assert.Equal(t, should, is)
	}
}

func refTime() time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	str1 := "2014-11-12T11:45:26.371Z"
	refTime, err := time.Parse(layout, str1)
	if err != nil {
		panic(err)
	}
	return refTime
}
