package reminderdaemon

import (
	"fmt"
	"sync"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/errors"
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
	i := 0
	for {
		i++
		start := time.Now()

		// Check for daily reminder every 5 minutes
		if i%5 == 0 {
			i = 0
			err := d.CheckForDailyReminder()
			if err != nil {
				log.Error(fmt.Sprintf("Error while checking daily reminders: %s", err.Error()))
			}
		}

		// Check for reminders every minute
		reminders, err := d.Database.GetPendingReminder()
		if err != nil {
			log.Warn("Not able to get Reminders from database: " + err.Error())
		} else {
			log.Info(fmt.Sprintf("REMINDERDAEMON: Found %d reminder to remind", len(reminders)))
			d.sendOutReminders(reminders)
		}

		sleepTime := time.Until(start.Add(time.Minute * 1))
		if sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}
	//wg.Done()
	//return nil
}

func (d *Daemon) sendOutReminders(reminders []database.Reminder) {
	for i := range reminders {
		originalMessage, err := d.Database.GetMessageFromReminder(reminders[i].ID, database.MessageTypeReminderRequest)
		if err != nil {
			log.Warn("Can not get original message: " + err.Error())
			continue
		}
		message, err := d.Messenger.SendReminder(&reminders[i], originalMessage)
		if err == errors.ErrEmptyChannel {
			_, err = d.Database.SetReminderDone(&reminders[i])
			if err != nil {
				log.Warn("Can not set reminder done: " + err.Error())
			}
			continue
		} else if err != nil {
			log.Warn(fmt.Sprintf("Failed to send reminder %d with: %s", reminders[i].ID, err.Error()))
			continue
		}

		_, err = d.Database.SetReminderDone(&reminders[i])
		if err != nil {
			log.Warn("Can not set reminder done: " + err.Error())
		}

		_, err = d.Database.AddMessage(message)
		if err != nil {
			log.Warn("Can not save message: " + err.Error())
		}
	}
}
