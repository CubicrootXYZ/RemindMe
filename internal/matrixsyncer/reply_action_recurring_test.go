package matrixsyncer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SyncerGetReplyActionRecurring(t *testing.T) {
	syncer := &Syncer{}
	action := syncer.getReplyActionRecurring(nil)

	assert.Greater(t, len(action.Name), 0)
	assert.Greater(t, len(action.Examples), 0)

	// Test the regex
	positiveCases := []string{"every 10 days   ", "every 99999 hours", "every hour "}
	positiveCases = append(positiveCases, action.Examples...)
	negativeCases := []string{" every 10 days"}

	for _, c := range positiveCases {
		assert.NotNil(t, action.Regex.Find([]byte(c)), "The following string should match: "+c)
	}

	for _, c := range negativeCases {
		assert.Nil(t, action.Regex.Find([]byte(c)), "The following string should not match: "+c)
	}
}
