package icalimporter

import (
	"testing"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	mocks "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/mocks"
	"github.com/tj/assert"

	"github.com/golang/mock/gomock"
)

func TestContentToReminders(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mocks.NewMockDatabase(ctrl)
	exampleTime, _ := time.Parse("2006-01-02T15:04:05", "2996-09-18T14:30:00")
	resource := testResource()

	db.EXPECT().AddOrUpdateThirdPartyResourceReminder(
		exampleTime,
		"A test event description",
		resource.Channel.ID,
		resource.ID,
		"uid1@example.com",
	)

	icalImporter := icalimporter{
		db: db,
	}

	err := icalImporter.contentToReminders(testEvent(), resource)

	assert.NoError(t, err)
}

func TestContentToReminders_RecurringEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := mocks.NewMockDatabase(ctrl)
	exampleTime, _ := time.Parse("2006-01-02T15:04:05", "2150-09-21T14:30:00")
	resource := testResource()

	db.EXPECT().AddOrUpdateThirdPartyResourceReminder(
		exampleTime,
		"A test event description",
		resource.Channel.ID,
		resource.ID,
		"uid1@example.com",
	)

	icalImporter := icalimporter{
		db: db,
	}

	err := icalImporter.contentToReminders(testEventRecurring(), resource)

	assert.NoError(t, err)
}

func testResource() *database.ThirdPartyResource {
	resource := &database.ThirdPartyResource{}
	resource.ID = 2
	resource.Channel.ID = 1

	return resource
}

func testEvent() string {
	return `BEGIN:VCALENDAR
PRODID:-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN
VERSION:2.0
BEGIN:VEVENT
DTSTAMP:19960704T120000Z
UID:uid1@example.com
ORGANIZER:mailto:jsmith@example.com
DTSTART:29960918T143000Z
DTEND:19960920T220000Z
STATUS:CONFIRMED
CATEGORIES:CONFERENCE
SUMMARY:A test event summary
DESCRIPTION:A test event description
END:VEVENT
END:VCALENDAR`
}

func testEventRecurring() string {
	return `BEGIN:VCALENDAR
PRODID:-//xyz Corp//NONSGML PDA Calendar Version 1.0//EN
VERSION:2.0
BEGIN:VEVENT
DTSTAMP:19960704T120000Z
UID:uid1@example.com
ORGANIZER:mailto:jsmith@example.com
DTSTART:21500918T143000Z
RRULE:FREQ=WEEKLY;BYDAY=MO
DTEND:21500918T153000Z
STATUS:CONFIRMED
CATEGORIES:CONFERENCE
SUMMARY:A test event summary
DESCRIPTION:A test event description
END:VEVENT
END:VCALENDAR`
}
