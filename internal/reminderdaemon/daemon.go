package reminderdaemon

import (
	"fmt"
	"sync"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
)

// Daemon holds all information for the reminder daemon
type Daemon struct {
	Database  Database
	Messenger Messenger
}

// Create creates a new reminder daemon
func Create(db Database, messenger Messenger) *Daemon {
	return &Daemon{
		Database:  db,
		Messenger: messenger,
	}
}

// Start starts the daemon
func (d *Daemon) Start(wg *sync.WaitGroup) error {
	for {
		start := time.Now()
		reminders, err := d.Database.GetPendingReminder()
		if err != nil {
			log.Warn("Not able to get Reminders from database: " + err.Error())
		} else {
			log.Info(fmt.Sprintf("REMINDERDAEMON: Found %d reminder to remind", len(*reminders)))
			for _, reminder := range *reminders {
				originalMessage, err := d.Database.GetMessageFromReminder(reminder.ID, database.MessageTypeReminderRequest)
				if err != nil {
					log.Warn("Can not get original message: " + err.Error())
					continue
				}
				message, err := d.Messenger.SendReminder(&reminder, originalMessage)
				if err != nil {
					log.Warn(fmt.Sprintf("Failed to send reminder %d with: %s", reminder.ID, err.Error()))
					continue
				}

				_, err = d.Database.SetReminderDone(&reminder)
				if err != nil {
					log.Warn("Can not set reminder done: " + err.Error())
				}

				_, err = d.Database.AddMessage(message)
				if err != nil {
					log.Warn("Can not save message: " + err.Error())
				}
			}
		}

		sleepTime := start.Add(time.Minute * 1).Sub(time.Now())
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}

	wg.Done()
	return nil
}
