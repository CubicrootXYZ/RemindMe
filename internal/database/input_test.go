package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInputByID(t *testing.T) {
	channel, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	inputBefore := testInput()

	err = service.AddInputToChannel(channel.ID, inputBefore)
	require.NoError(t, err)

	inputAfter, err := service.GetInputByID(inputBefore.ID)
	require.NoError(t, err)

	assert.Equal(t, inputBefore.ID, inputAfter.ID)
	assert.Equal(t, inputBefore.InputID, inputAfter.InputID)
	assert.Equal(t, inputBefore.InputType, inputAfter.InputType)
}

func TestGetInputByType(t *testing.T) {
	channel, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	inputBefore := testInput()

	err = service.AddInputToChannel(channel.ID, inputBefore)
	require.NoError(t, err)

	inputAfter, err := service.GetInputByType(inputBefore.InputID, inputBefore.InputType)
	require.NoError(t, err)

	assert.Equal(t, inputBefore.ID, inputAfter.ID)
	assert.Equal(t, inputBefore.InputID, inputAfter.InputID)
	assert.Equal(t, inputBefore.InputType, inputAfter.InputType)
}
