package reminderdaemon

import (
	"fmt"
	"sync"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/asyncmessenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
)

// Daemon holds all information for the reminder daemon
type Daemon struct {
	Database  Database
	Messenger asyncmessenger.Messenger
}

// Create creates a new reminder daemon
func Create(db Database, messenger asyncmessenger.Messenger) *Daemon {
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

		go d.sendReminder(&reminders[i], originalMessage)
	}
}

func (d *Daemon) sendReminder(reminder *database.Reminder, originalMessage *database.Message) {
	if reminder.Channel.ID == 0 {
		log.Error("Can not send reminders to empty channels")
		return
	}

	message, messageFormatted := formater.ReminderToMessage(reminder)

	responseMessage := &asyncmessenger.Response{
		Message:                   message,
		MessageFormatted:          messageFormatted,
		RespondToMessage:          originalMessage.Body,
		RespondToMessageFormatted: originalMessage.BodyHTML,
		RespondToEventID:          originalMessage.ExternalIdentifier,
		ChannelExternalIdentifier: reminder.Channel.ChannelIdentifier,
	}
	resp, err := d.Messenger.SendResponse(responseMessage)
	if err != nil {
		log.Error("Failed to send reminder message: " + err.Error())
		return
	}

	for _, reaction := range types.ReactionsReminder {
		err = d.Messenger.SendReactionAsync(&asyncmessenger.Reaction{
			Reaction:                  reaction,
			MessageExternalIdentifier: resp.ExternalIdentifier,
			ChannelExternalIdentifier: reminder.Channel.ChannelIdentifier,
		})
		if err != nil {
			log.Error("Failed to send reaction " + reaction + " to reminder message: " + err.Error())
		}
	}
	_, err = d.Database.AddMessage(&database.Message{
		Body:               message,
		BodyHTML:           messageFormatted,
		ReminderID:         &reminder.ID,
		ChannelID:          reminder.ChannelID,
		Type:               database.MessageTypeReminder,
		Timestamp:          resp.Timestamp,
		ExternalIdentifier: resp.ExternalIdentifier,
	})
	if err != nil {
		log.Error("Failed saving reminder message: " + err.Error())
	}

	_, err = d.Database.SetReminderDone(reminder)
	if err != nil {
		log.Error("Failed setting reminder done: " + err.Error())
	}
}
