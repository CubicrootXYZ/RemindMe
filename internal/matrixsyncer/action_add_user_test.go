package matrixsyncer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SyncerGetActionAddUser(t *testing.T) {
	syncer := &Syncer{}
	action := syncer.getActionAddUser()

	assert.Greater(t, len(action.Name), 0)
	assert.Greater(t, len(action.Examples), 0)

	// Test the regex
	positiveCases := []string{"add user a", "add user ", "add user @user:instance.org", "add user abcdefghijkl mandsjfue"}
	positiveCases = append(positiveCases, action.Examples...)
	negativeCases := []string{"may I add user a?", " add user a", "add no user"}

	for _, c := range positiveCases {
		assert.NotNil(t, action.Regex.Find([]byte(c)), "The following string should match: "+c)
	}

	for _, c := range negativeCases {
		assert.Nil(t, action.Regex.Find([]byte(c)), "The following string should not match: "+c)
	}
}
