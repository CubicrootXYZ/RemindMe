package matrixsyncer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SyncerGetReplyActionDelete(t *testing.T) {
	syncer := &Syncer{}
	action := syncer.getReplyActionDelete(nil)

	assert.Greater(t, len(action.Name), 0)
	assert.Greater(t, len(action.Examples), 0)

	// Test the regex
	positiveCases := []string{"delete   ", "remove "}
	positiveCases = append(positiveCases, action.Examples...)
	negativeCases := []string{" delete", " delete reminder", "delete this"}

	for _, c := range positiveCases {
		assert.NotNil(t, action.Regex.Find([]byte(c)), "The following string should match: "+c)
	}

	for _, c := range negativeCases {
		assert.Nil(t, action.Regex.Find([]byte(c)), "The following string should not match: "+c)
	}
}
