package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOutputByID(t *testing.T) {
	channel, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	outputBefore := testOutput()

	err = service.AddOutputToChannel(channel.ID, outputBefore)
	require.NoError(t, err)

	outputAfter, err := service.GetOutputByID(outputBefore.ID)
	require.NoError(t, err)

	assert.Equal(t, outputBefore.ID, outputAfter.ID)
	assert.Equal(t, outputBefore.OutputID, outputAfter.OutputID)
	assert.Equal(t, outputBefore.OutputType, outputAfter.OutputType)
	assert.Equal(t, outputBefore.LastDailyReminder, outputAfter.LastDailyReminder)
}

func TestUpdateOutput(t *testing.T) {
	channel, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	outputBefore := testOutput()

	err = service.AddOutputToChannel(channel.ID, outputBefore)
	require.NoError(t, err)

	outputBefore.OutputID = 123
	outputBefore.OutputType = "shrug"

	_, err = service.UpdateOutput(outputBefore)
	require.NoError(t, err)

	outputAfter, err := service.GetOutputByID(outputBefore.ID)
	require.NoError(t, err)

	assert.Equal(t, outputBefore.ID, outputAfter.ID)
	assert.Equal(t, outputBefore.OutputID, outputAfter.OutputID)
	assert.Equal(t, outputBefore.OutputType, outputAfter.OutputType)
	assert.Equal(t, outputBefore.LastDailyReminder, outputAfter.LastDailyReminder)
}
