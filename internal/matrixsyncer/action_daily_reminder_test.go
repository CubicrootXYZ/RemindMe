package matrixsyncer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SyncerGetActionSetDailyReminder(t *testing.T) {
	syncer := &Syncer{}
	action := syncer.getActionSetDailyReminder()

	assert.Greater(t, len(action.Name), 0)
	assert.Greater(t, len(action.Examples), 0)

	// Test the regex
	positiveCases := []string{"set   a   daily message at 8 am  ", "update the  daily reminder to 10:00   "}
	positiveCases = append(positiveCases, action.Examples...)
	negativeCases := []string{" set a daily message at 8 am"}

	for _, c := range positiveCases {
		assert.NotNil(t, action.Regex.Find([]byte(c)), "The following string should match: "+c)
	}

	for _, c := range negativeCases {
		assert.Nil(t, action.Regex.Find([]byte(c)), "The following string should not match: "+c)
	}
}
