package daemon

import "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"

func outputFromDatabase(output *database.Output) *Output {
	return &Output{
		ID:         output.ID,
		OutputType: output.OutputType,
		OutputID:   output.OutputID,
	}
}
