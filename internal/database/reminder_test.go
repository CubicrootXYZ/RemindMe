package database

import "github.com/DATA-DOG/go-sqlmock"

func rowsForReminders(reminders []*Reminder) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "remind_time", "message", "active", "repeat_interval", "repeat_max", "repeated", "channel_id"})

	for _, reminder := range reminders {
		rows.AddRow(
			reminder.ID,
			reminder.CreatedAt,
			reminder.UpdatedAt,
			reminder.DeletedAt,
			reminder.RemindTime,
			reminder.Message,
			reminder.Active,
			reminder.RepeatInterval,
			reminder.RepeatMax,
			reminder.Repeated,
			reminder.ChannelID,
		)
	}

	return rows
}
