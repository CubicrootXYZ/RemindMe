package formater

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormater_ToNiceDuration(t *testing.T) {
	testCases := make(map[time.Duration]string)
	testCases[5*time.Millisecond] = "0 seconds"
	testCases[time.Second] = "1 seconds"
	testCases[59*time.Second] = "59 seconds"
	testCases[99*time.Second] = "1 minutes"
	testCases[5*time.Minute] = "5 minutes"
	testCases[24*time.Minute] = "24 minutes"
	testCases[24*time.Minute+500*time.Microsecond] = "24 minutes"
	testCases[time.Hour] = "1 hours"
	testCases[36*time.Hour+59*time.Minute] = "36 hours"
	testCases[48*time.Hour+1*time.Microsecond] = "2 days"
	testCases[46879*24*time.Hour] = "46879 days"

	for duration, should := range testCases {
		is := ToNiceDuration(duration)
		assert.Equal(t, should, is)
	}
}
