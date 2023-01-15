package daemon

import "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"

type service struct {
	OutputServices map[string]OutputService // Maps OutputTypes to the services
	Database       database.Service
}

// NewService assembles a new service.
func NewService() Service {
	return service{}
}
