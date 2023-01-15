package database_test

import (
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testChannel() *database.Channel {
	dailyReminder := uint(120)
	return &database.Channel{
		Description:   "channel description",
		DailyReminder: &dailyReminder,
	}
}

func testInput() *database.Input {
	return &database.Input{
		InputType: "test",
		InputID:   1,
		Enabled:   true,
	}
}

func testOutput() *database.Output {
	return &database.Output{
		OutputType: "test",
		OutputID:   1,
		Enabled:    true,
	}
}

func TestService_NewChannel(t *testing.T) {
	channelBefore := testChannel()

	channelAfter, err := service.NewChannel(channelBefore)
	require.NoError(t, err)

	assert.Equal(t, channelBefore.Description, channelAfter.Description)
	assert.Equal(t, channelBefore.DailyReminder, channelAfter.DailyReminder)
	assert.Less(t, uint(0), channelAfter.ID)
	assert.False(t, channelAfter.CreatedAt.IsZero())
	assert.False(t, channelAfter.UpdatedAt.IsZero())
}

func TestService_GetChannelByID(t *testing.T) {
	channelBefore, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	assert.Equal(t, channelBefore.ID, channelAfter.ID)
	assert.Equal(t, channelBefore.Description, channelAfter.Description)
	assert.Equal(t, channelBefore.DailyReminder, channelAfter.DailyReminder)
	assert.Less(t, uint(0), channelAfter.ID)
	assert.False(t, channelAfter.CreatedAt.IsZero())
	assert.False(t, channelAfter.UpdatedAt.IsZero())
}

func TestService_GetChannelByIDWithNotFound(t *testing.T) {
	_, err := service.GetChannelByID(123456)
	require.Error(t, err)

	assert.ErrorIs(t, err, database.ErrNotFound)
}

func TestService_AddInputToChannel(t *testing.T) {
	inputBefore := testInput()

	channelBefore, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	err = service.AddInputToChannel(channelBefore.ID, inputBefore)
	require.NoError(t, err)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	require.Equal(t, 1, len(channelAfter.Inputs))
	assert.Equal(t, inputBefore.InputType, channelAfter.Inputs[0].InputType)
	assert.Equal(t, inputBefore.InputID, channelAfter.Inputs[0].InputID)
	assert.Equal(t, inputBefore.Enabled, channelAfter.Inputs[0].Enabled)
}

func TestService_AddOutputToChannel(t *testing.T) {
	outputBefore := testOutput()

	channelBefore, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	err = service.AddOutputToChannel(channelBefore.ID, outputBefore)
	require.NoError(t, err)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	require.Equal(t, 1, len(channelAfter.Outputs))
	assert.Equal(t, outputBefore.OutputType, channelAfter.Outputs[0].OutputType)
	assert.Equal(t, outputBefore.OutputID, channelAfter.Outputs[0].OutputID)
	assert.Equal(t, outputBefore.Enabled, channelAfter.Outputs[0].Enabled)
}
