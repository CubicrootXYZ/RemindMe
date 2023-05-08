package database_test

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testOutput() *database.IcalOutput {
	return &database.IcalOutput{
		Token: "abcde",
	}
}

func TestService_NewIcalOutput(t *testing.T) {
	start := time.Now()
	time.Sleep(time.Millisecond) // Avoids issues with database time representation being less accurate.
	outputBefore, err := service.NewIcalOutput(testOutput())
	require.NoError(t, err)

	assert.Greater(t, outputBefore.ID, uint(0))
	assert.GreaterOrEqual(t, outputBefore.CreatedAt, start)
	assert.NotEmpty(t, outputBefore.Token)

	outputAfter, err := service.GetIcalOutputByID(outputBefore.ID)
	require.NoError(t, err)
	assert.Equal(t, outputBefore.ID, outputAfter.ID)
	assert.Equal(t, outputBefore.Token, outputAfter.Token)
}

func TestService_DeleteIcalOutput(t *testing.T) {
	output, err := service.NewIcalOutput(testOutput())
	require.NoError(t, err)
	require.Greater(t, output.ID, uint(0))

	err = service.DeleteIcalOutput(output.ID)
	require.NoError(t, err)

	_, err = service.GetIcalOutputByID(output.ID)
	require.ErrorIs(t, err, database.ErrNotFound)
}

func TestService_DeleteIcalOutputWithNotFound(t *testing.T) {
	err := service.DeleteIcalOutput(999999)
	require.ErrorIs(t, err, database.ErrNotFound)
}

func TestService_GenerateNewToken(t *testing.T) {
	outputBefore, err := service.NewIcalOutput(testOutput())
	require.NoError(t, err)

	output := *outputBefore
	outputAfter, err := service.GenerateNewToken(&output)
	require.NoError(t, err)

	assert.NotEqual(t, outputBefore.Token, outputAfter.Token)
	assert.Greater(t, len(outputAfter.Token), 29)
}
