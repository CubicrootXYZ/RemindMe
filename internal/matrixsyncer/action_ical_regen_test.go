package matrixsyncer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SyncerGetActionIcalRegenerate(t *testing.T) {
	syncer := &Syncer{}
	action := syncer.getActionIcalRegenerate()

	assert.Greater(t, len(action.Name), 0)
	assert.Greater(t, len(action.Examples), 0)

	// Test the regex
	positiveCases := []string{"make  renew    the    ical ", "generate new secret  "}
	positiveCases = append(positiveCases, action.Examples...)
	negativeCases := []string{" generate new secret", "generate new secret key token"}

	for _, c := range positiveCases {
		assert.NotNil(t, action.Regex.Find([]byte(c)), "The following string should match: "+c)
	}

	for _, c := range negativeCases {
		assert.Nil(t, action.Regex.Find([]byte(c)), "The following string should not match: "+c)
	}
}
