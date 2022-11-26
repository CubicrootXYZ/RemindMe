package matrixsyncer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SyncerGetActionChangeReminder(t *testing.T) {
	syncer := &Syncer{}
	action := syncer.getActionChangeReminder()

	assert.Greater(t, len(action.Name), 0)
	assert.Greater(t, len(action.Examples), 0)

	// Test the regex
	positiveCases := []string{"change reminder 1 to tomorrow", "update 68 to Saturday 4 pm", "set reminder 1 to Saturday 16:00", "update reminder id 19488382 to next week", "update  reminder id 1   to tomorrow", "update reminder id1 to tomorrow", "update reminder 1 to tomorrow    ", "update 1 to tomorrow"}
	positiveCases = append(positiveCases, action.Examples...)
	negativeCases := []string{" change reminder 1 to tomorrow", "set reminder a to tomorrow"}

	for _, c := range positiveCases {
		assert.NotNil(t, action.Regex.Find([]byte(c)), "The following string should match: "+c)
	}

	for _, c := range negativeCases {
		assert.Nil(t, action.Regex.Find([]byte(c)), "The following string should not match: "+c)
	}
}
