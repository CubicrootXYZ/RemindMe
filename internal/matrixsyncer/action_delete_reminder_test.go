package matrixsyncer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SyncerGetActionDeleteReminder(t *testing.T) {
	syncer := &Syncer{}
	action := syncer.getActionDeleteReminder()

	assert.Greater(t, len(action.Name), 0)
	assert.Greater(t, len(action.Examples), 0)

	// Test the regex
	positiveCases := []string{"delete reminder 1   ", "delete   reminder   1", "remove reminder 118485832  ", "delete  1  ", "remove 99999  "}
	positiveCases = append(positiveCases, action.Examples...)
	negativeCases := []string{" delete reminder 1", "delete1", " delete 1"}

	for _, c := range positiveCases {
		assert.NotNil(t, action.Regex.Find([]byte(c)), "The following string should match: "+c)
	}

	for _, c := range negativeCases {
		assert.Nil(t, action.Regex.Find([]byte(c)), "The following string should not match: "+c)
	}
}
