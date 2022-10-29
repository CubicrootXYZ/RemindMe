package icalimporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"

	ical "github.com/arran4/golang-ical"
)

type icalimporter struct {
	stop          chan bool
	db            types.Database
	reminderDelay time.Duration // Time difference the reminder should have compored to the ical event
}

func NewIcalImporter(db types.Database) IcalImporter {
	return &icalimporter{
		stop:          make(chan bool),
		reminderDelay: time.Minute * -5,
		db:            db,
	}
}

// Run runs the importer, call this within a goroutine.
func (importer *icalimporter) Run() {
	ticker := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-importer.stop:
			log.Info("Stopping importer")
			return
		case <-ticker.C:
			importer.updateIcalResources()
		}
	}
}

// Stop shuts the importer down nicely.
func (importer *icalimporter) Stop() {
	importer.stop <- true
}

func (importer *icalimporter) updateIcalResources() {
	icalResources, err := importer.db.GetThirdPartyResources(database.ThirdPartyResourceTypeIcal)
	if err != nil {
		log.Error("Can not fetch third party resources from database: " + err.Error())
	}

	for i := range icalResources {
		if icalResources[i].Channel.ID == 0 {
			log.Error("Can not fetch resources for empty channels")
			continue
		}

		content, err := getFileContent(icalResources[i].ResourceURL)
		if err != nil {
			log.Error(fmt.Sprintf("Failed fetching the resource %s for channel ID %d: %s", icalResources[i].ResourceURL, icalResources[i].ChannelID, err.Error()))
			continue
		}

		err = importer.contentToReminders(content, &icalResources[i])
		if err != nil {
			log.Error("Failed to parse third party resources content to a reminder: " + err.Error())
			continue
		}
	}
}

func (importer *icalimporter) contentToReminders(content string, resource *database.ThirdPartyResource) error {
	calendar, err := ical.ParseCalendar(strings.NewReader(content))
	if err != nil {
		return err
	}

	for _, event := range calendar.Events() {
		idProp := event.GetProperty(ical.ComponentPropertyUniqueId)
		if idProp == nil {
			log.Info("Skipping event, can not read id")
			continue
		}
		id := idProp.Value
		if len(id) <= 2 {
			log.Info("Skipping event, id is to short")
			continue
		}

		startTime, err := event.GetStartAt()
		if err != nil {
			log.Info("Skipping event, can not read start time: " + err.Error())
			continue
		}
		startTime = startTime.Add(importer.reminderDelay)

		if time.Until(startTime) < 0 {
			// Ignore past events
			continue
		}

		name := getNameFromEvent(event)
		if name == "" {
			log.Info("Skipping event, can not read name")
			continue
		}

		_, err = importer.db.AddOrUpdateThirdPartyResourceReminder(startTime, name, resource.Channel.ID, resource.ID, id)
		if err != nil {
			return err
		}
	}

	return nil
}
