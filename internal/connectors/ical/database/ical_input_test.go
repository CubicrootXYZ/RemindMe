package database_test

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInput() *database.IcalInput {
	return &database.IcalInput{
		URL: "https://example.com",
	}
}

func TestService_NewIcalInput(t *testing.T) {
	start := time.Now()
	time.Sleep(time.Millisecond) // Avoids issues with database time representation being less accurate.
	inputBefore, err := service.NewIcalInput(testInput())
	require.NoError(t, err)

	assert.Greater(t, inputBefore.ID, uint(0))
	assert.GreaterOrEqual(t, inputBefore.CreatedAt, start)
	assert.NotEmpty(t, inputBefore.URL)

	inputAfter, err := service.GetIcalInputByID(inputBefore.ID)
	require.NoError(t, err)
	assert.Equal(t, inputBefore.ID, inputAfter.ID)
	assert.Equal(t, inputBefore.URL, inputAfter.URL)
}

func TestService_DeleteIcalInput(t *testing.T) {
	input, err := service.NewIcalInput(testInput())
	require.NoError(t, err)
	require.Greater(t, input.ID, uint(0))

	err = service.DeleteIcalInput(input.ID)
	require.NoError(t, err)

	_, err = service.GetIcalInputByID(input.ID)
	require.ErrorIs(t, err, database.ErrNotFound)
}

func TestService_DeleteIcalInputWithNotFound(t *testing.T) {
	err := service.DeleteIcalInput(999999)
	require.ErrorIs(t, err, database.ErrNotFound)
}

func TestService_ListIcalInput(t *testing.T) {
	input, err := service.NewIcalInput(testInput())
	require.NoError(t, err)
	require.Greater(t, input.ID, uint(0))

	inputs, err := service.ListIcalInputs(&database.ListIcalInputsOpts{})
	require.NoError(t, err)

	found := false
	for _, i := range inputs {
		if i.ID == input.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "entity is not in list return")

	tr := true
	inputs, err = service.ListIcalInputs(&database.ListIcalInputsOpts{
		Disabled: &tr,
	})
	require.NoError(t, err)

	found = false
	for _, i := range inputs {
		if i.ID == input.ID {
			found = true
			break
		}
	}
	assert.False(t, found, "entity is in list return")
}

func TestService_UpdateIcalInput(t *testing.T) {
	input, err := service.NewIcalInput(testInput())
	require.NoError(t, err)
	require.Greater(t, input.ID, uint(0))

	input.URL = "abcdef"

	inputAfter, err := service.UpdateIcalInput(input)
	require.NoError(t, err)

	assert.Equal(t, input, inputAfter)
}
