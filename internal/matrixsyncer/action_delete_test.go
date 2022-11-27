package matrixsyncer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SyncerGetActionDelete(t *testing.T) {
	syncer := &Syncer{}
	action := syncer.getActionDelete()

	assert.Greater(t, len(action.Name), 0)
	assert.Greater(t, len(action.Examples), 0)

	// Test the regex
	positiveCases := []string{"delete the data from remindme", "delete the data at remindme   "}
	positiveCases = append(positiveCases, action.Examples...)
	negativeCases := []string{" delete the data from remindme", " delete", "delete data remindme"}

	for _, c := range positiveCases {
		assert.NotNil(t, action.Regex.Find([]byte(c)), "The following string should match: "+c)
	}

	for _, c := range negativeCases {
		assert.Nil(t, action.Regex.Find([]byte(c)), "The following string should not match: "+c)
	}
}
