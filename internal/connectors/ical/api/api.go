package api

import "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/ical/database"

type api struct {
	icalDB database.Service
}
