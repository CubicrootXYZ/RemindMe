package database_test

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database/mocks"
	"github.com/golang/mock/gomock"
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
	id := uint(1)
	_, err := service.GetInputByID(id)
	for !errors.Is(err, database.ErrNotFound) {
		id = uint(rand.Int()) //nolint:gosec
		_, err = service.GetInputByType(id, "test")
	}

	return &database.Input{
		InputType: "test",
		InputID:   id,
		Enabled:   true,
	}
}

func testOutput() *database.Output {
	id := uint(1)
	_, err := service.GetOutputByID(id)
	for !errors.Is(err, database.ErrNotFound) {
		id = uint(rand.Int()) //nolint:gosec
		_, err = service.GetOutputByType(id, "test")
	}

	return &database.Output{
		OutputType: "test",
		OutputID:   id,
		Enabled:    true,
	}
}

func TestService_NewChannel(t *testing.T) {
	channelBefore := testChannel()

	channelAfter, err := service.NewChannel(channelBefore)
	require.NoError(t, err)

	assert.Equal(t, channelBefore.Description, channelAfter.Description)
	assert.Equal(t, channelBefore.DailyReminder, channelAfter.DailyReminder)
	assert.NotZero(t, channelAfter.ID)
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
	assert.NotZero(t, channelAfter.ID)
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

	require.Len(t, channelAfter.Inputs, 1)
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

	require.Len(t, channelAfter.Outputs, 1)
	assert.Equal(t, outputBefore.OutputType, channelAfter.Outputs[0].OutputType)
	assert.Equal(t, outputBefore.OutputID, channelAfter.Outputs[0].OutputID)
	assert.Equal(t, outputBefore.Enabled, channelAfter.Outputs[0].Enabled)
}

func TestService_RemoveInputFromChannel(t *testing.T) { //nolint:dupl // wrong hint
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	input := testInput()
	inputService := mocks.NewMockInputService(ctrl)
	inputService.EXPECT().InputRemoved(input.InputType, input.InputID).Return(nil)

	service := getService(&database.Config{
		Connection:    getConnection(),
		InputServices: map[string]database.InputService{"test": inputService},
	})

	channelBefore, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	err = service.AddInputToChannel(channelBefore.ID, input)
	require.NoError(t, err)

	err = service.RemoveInputFromChannel(channelBefore.ID, input.ID)
	require.NoError(t, err)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	assert.Empty(t, channelAfter.Inputs)

	_, err = service.GetInputByID(input.ID)
	require.Error(t, err)
	assert.ErrorIs(t, err, database.ErrNotFound)
}

func TestService_RemoveInputFromChannelWithInputServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	input := testInput()
	expectedErr := errors.New("test")

	inputService := mocks.NewMockInputService(ctrl)
	inputService.EXPECT().InputRemoved(input.InputType, input.InputID).Return(expectedErr)

	service := getService(&database.Config{
		Connection:    getConnection(),
		InputServices: map[string]database.InputService{"test": inputService},
	})

	channelBefore, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	err = service.AddInputToChannel(channelBefore.ID, input)
	require.NoError(t, err)

	err = service.RemoveInputFromChannel(channelBefore.ID, input.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	require.Len(t, channelAfter.Inputs, 1)
}

func TestService_RemoveInputFromChannelWithoutInputService(t *testing.T) {
	input := testInput()

	channelBefore, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	err = service.AddInputToChannel(channelBefore.ID, input)
	require.NoError(t, err)

	err = service.RemoveInputFromChannel(channelBefore.ID, input.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, database.ErrUnknownInput)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	require.Len(t, channelAfter.Inputs, 1)
}

func TestService_RemoveOutputFromChannel(t *testing.T) { //nolint:dupl // wrong hint
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	output := testOutput()
	outputService := mocks.NewMockOutputService(ctrl)
	outputService.EXPECT().OutputRemoved(output.OutputType, output.OutputID).Return(nil)

	service := getService(&database.Config{
		Connection:     getConnection(),
		OutputServices: map[string]database.OutputService{"test": outputService},
	})

	channelBefore, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	err = service.AddOutputToChannel(channelBefore.ID, output)
	require.NoError(t, err)

	err = service.RemoveOutputFromChannel(channelBefore.ID, output.ID)
	require.NoError(t, err)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	assert.Empty(t, channelAfter.Outputs)

	_, err = service.GetOutputByID(output.ID)
	require.Error(t, err)
	assert.ErrorIs(t, err, database.ErrNotFound)
}

func TestService_RemoveOutputFromChannelWithOutputServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	output := testOutput()
	expectedErr := errors.New("test")

	outputService := mocks.NewMockOutputService(ctrl)
	outputService.EXPECT().OutputRemoved(output.OutputType, output.OutputID).Return(expectedErr)

	service := getService(&database.Config{
		Connection:     getConnection(),
		OutputServices: map[string]database.OutputService{"test": outputService},
	})

	channelBefore, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	err = service.AddOutputToChannel(channelBefore.ID, output)
	require.NoError(t, err)

	err = service.RemoveOutputFromChannel(channelBefore.ID, output.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	require.Len(t, channelAfter.Outputs, 1)
}

func TestService_RemoveOutputFromChannelWithoutOutputService(t *testing.T) {
	output := testOutput()

	channelBefore, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	err = service.AddOutputToChannel(channelBefore.ID, output)
	require.NoError(t, err)

	err = service.RemoveOutputFromChannel(channelBefore.ID, output.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, database.ErrUnknownOutput)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	assert.Len(t, channelAfter.Outputs, 1)
}

func TestService_UpdateChannel(t *testing.T) {
	channelBefore := testChannel()

	channelBefore, err := service.NewChannel(channelBefore)
	require.NoError(t, err)

	dailyReminder := uint(120)
	channelBefore.Description = "test 2"
	channelBefore.DailyReminder = &dailyReminder

	_, err = service.UpdateChannel(channelBefore)
	require.NoError(t, err)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	assert.Equal(t, channelBefore.Description, channelAfter.Description)
	assert.Equal(t, channelBefore.DailyReminder, channelAfter.DailyReminder)
}

func TestService_GetChannels(t *testing.T) {
	channelBefore, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	channels, err := service.GetChannels()
	require.NoError(t, err)

	channelFound := false
	for _, channelAfter := range channels {
		if channelAfter.ID == channelBefore.ID {
			channelFound = true
		}
	}

	assert.True(t, channelFound)
}

func TestService_DeleteChannel(t *testing.T) {
	channel, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	err = service.DeleteChannel(channel.ID)
	assert.NoError(t, err)
}

func TestService_DeleteChannelWithChannelNotFound(t *testing.T) {
	err := service.DeleteChannel(999999999)
	assert.ErrorIs(t, err, database.ErrNotFound)
}
