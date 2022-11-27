package matrixsyncer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SyncerGetActionList(t *testing.T) {
	syncer := &Syncer{}
	action := syncer.getActionList()

	assert.Greater(t, len(action.Name), 0)
	assert.Greater(t, len(action.Examples), 0)

	// Test the regex
	positiveCases := []string{"list all reminders please   ", "show my reminders   "}
	positiveCases = append(positiveCases, action.Examples...)
	negativeCases := []string{" show my reminders"}

	for _, c := range positiveCases {
		assert.NotNil(t, action.Regex.Find([]byte(c)), "The following string should match: "+c)
	}

	for _, c := range negativeCases {
		assert.Nil(t, action.Regex.Find([]byte(c)), "The following string should not match: "+c)
	}
}
