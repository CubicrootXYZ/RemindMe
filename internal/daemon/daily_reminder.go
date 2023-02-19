package daemon

import "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"

func dailyReminderFromDatabase(events []database.Event) *DailyReminder {
	dailyReminder := &DailyReminder{
		Events: make([]Event, len(events)),
	}

	for i := range events {
		dailyReminder.Events[i] = *eventFromDatabase(&events[i])
	}

	return dailyReminder
}
