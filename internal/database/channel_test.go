package database_test

import (
	"errors"
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

func TestRemoveInputFromChannel(t *testing.T) { //nolint:dupl // wrong hint
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	input := testInput()
	inputService := mocks.NewMockInputService(ctrl)
	inputService.EXPECT().InputRemoved(input.InputType, input.InputID, gomock.Any()).Return(nil)

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

	assert.Equal(t, 0, len(channelAfter.Inputs))

	_, err = service.GetInputByID(input.ID)
	require.Error(t, err)
	assert.ErrorIs(t, err, database.ErrNotFound)
}

func TestRemoveInputFromChannelWithInputServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	input := testInput()
	expectedErr := errors.New("test")

	inputService := mocks.NewMockInputService(ctrl)
	inputService.EXPECT().InputRemoved(input.InputType, input.InputID, gomock.Any()).Return(expectedErr)

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
	assert.ErrorIs(t, err, expectedErr)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	assert.Equal(t, 1, len(channelAfter.Inputs))
}

func TestRemoveInputFromChannelWithoutInputService(t *testing.T) {
	input := testInput()

	channelBefore, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	err = service.AddInputToChannel(channelBefore.ID, input)
	require.NoError(t, err)

	err = service.RemoveInputFromChannel(channelBefore.ID, input.ID)
	require.Error(t, err)
	assert.ErrorIs(t, err, database.ErrUnknownInput)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	assert.Equal(t, 1, len(channelAfter.Inputs))
}

func TestRemoveOutputFromChannel(t *testing.T) { //nolint:dupl // wrong hint
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	output := testOutput()
	outputService := mocks.NewMockOutputService(ctrl)
	outputService.EXPECT().OutputRemoved(output.OutputType, output.OutputID, gomock.Any()).Return(nil)

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

	assert.Equal(t, 0, len(channelAfter.Outputs))

	_, err = service.GetOutputByID(output.ID)
	require.Error(t, err)
	assert.ErrorIs(t, err, database.ErrNotFound)
}

func TestRemoveOutputFromChannelWithOutputServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	output := testOutput()
	expectedErr := errors.New("test")

	outputService := mocks.NewMockOutputService(ctrl)
	outputService.EXPECT().OutputRemoved(output.OutputType, output.OutputID, gomock.Any()).Return(expectedErr)

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
	assert.ErrorIs(t, err, expectedErr)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	assert.Equal(t, 1, len(channelAfter.Outputs))
}

func TestRemoveOutputFromChannelWithoutOutputService(t *testing.T) {
	output := testOutput()

	channelBefore, err := service.NewChannel(testChannel())
	require.NoError(t, err)

	err = service.AddOutputToChannel(channelBefore.ID, output)
	require.NoError(t, err)

	err = service.RemoveOutputFromChannel(channelBefore.ID, output.ID)
	require.Error(t, err)
	assert.ErrorIs(t, err, database.ErrUnknownOutput)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	assert.Equal(t, 1, len(channelAfter.Outputs))
}

func TestUpdateChannel(t *testing.T) {
	channelBefore := testChannel()

	channelBefore, err := service.NewChannel(channelBefore)
	require.NoError(t, err)

	dailyReminder := uint(120)
	lastDailyReminder := time2123().UTC()
	channelBefore.Description = "test 2"
	channelBefore.DailyReminder = &dailyReminder
	channelBefore.LastDailyReminder = &lastDailyReminder

	_, err = service.UpdateChannel(channelBefore)
	require.NoError(t, err)

	channelAfter, err := service.GetChannelByID(channelBefore.ID)
	require.NoError(t, err)

	assert.Equal(t, channelBefore.Description, channelAfter.Description)
	assert.Equal(t, channelBefore.DailyReminder, channelAfter.DailyReminder)
	assert.Equal(t, channelBefore.LastDailyReminder, channelAfter.LastDailyReminder)
}
