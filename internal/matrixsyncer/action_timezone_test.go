package matrixsyncer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SyncerGetActionTimezone(t *testing.T) {
	syncer := &Syncer{}
	action := syncer.getActionTimezone()

	assert.Greater(t, len(action.Name), 0)
	assert.Greater(t, len(action.Examples), 0)

	// Test the regex
	positiveCases := []string{"set timezone to jdhgkldfhglkdf djfhgk dfg   "}
	positiveCases = append(positiveCases, action.Examples...)
	negativeCases := []string{" set timezone to Europe/Berlin", " set timezone", "timezone"}

	for _, c := range positiveCases {
		assert.NotNil(t, action.Regex.Find([]byte(c)), "The following string should match: "+c)
	}

	for _, c := range negativeCases {
		assert.Nil(t, action.Regex.Find([]byte(c)), "The following string should not match: "+c)
	}
}
