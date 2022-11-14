package formater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetSuffixIntOnSuccess(t *testing.T) {
	result, err := GetSuffixInt("Well, this is a very long sentencen: please add some more characters and an integer like 203")
	assert.NoError(t, err)
	assert.Equal(t, 203, result)
}

func Test_GetSuffixIntOnFailure(t *testing.T) {
	result, err := GetSuffixInt("Well, this is a very long sentencen: please add some more characters and an integer like two hundred")
	assert.Error(t, err)
	assert.Equal(t, 0, result)
}
