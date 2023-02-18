package message_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/actions/message"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewEventAction(t *testing.T) {
	action := &message.NewEventAction{}

	assert.NotEmpty(t, action.Name())

	title, desc, examples := action.GetDocu()
	assert.NotEmpty(t, title)
	assert.NotEmpty(t, desc)
	assert.NotEmpty(t, examples)

	assert.NotNil(t, action.Selector())
}

func TestNewEventAction_Selector(t *testing.T) {
	action := &message.NewEventAction{}

	r := action.Selector()
	assert.NotNil(t, r)
}

func TestNewEventAction_HandleEvent(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := database.NewMockService(ctrl)
	matrixDB := matrixdb.NewMockService(ctrl)
	client := mautrixcl.NewMockClient(ctrl)
	msngr := messenger.NewMockMessenger(ctrl)

	action := &message.NewEventAction{}
	action.Configure(
		gologger.New(gologger.LogLevelDebug, 0),
		client,
		msngr,
		matrixDB,
		db,
	)

	msgs := []string{
		"my test reminder at monday 1:11",
		"my +#ä§$&7(&$==§ é reminder in 100 years",
	}

	for _, msg := range msgs {
		t.Run(msg, func(t *testing.T) {
			// Expectations
			db.EXPECT().NewEvent(&eventMatcher{
				evt: &database.Event{
					Duration:  message.DefaultEventTime,
					Message:   msg,
					Active:    true,
					ChannelID: testEvent().Channel.ID,
					InputID:   &testEvent().Input.ID,
				},
			}).Return(&database.Event{
				Model: gorm.Model{
					ID: 1,
				},
				Duration:  message.DefaultEventTime,
				Message:   msg,
				Active:    true,
				ChannelID: testEvent().Channel.ID,
				InputID:   &testEvent().Input.ID,
			}, nil)

			matrixDB.EXPECT().NewMessage(&matrixdb.MatrixMessage{
				UserID:        testEvent().Event.Sender.String(),
				RoomID:        testEvent().Room.ID,
				Body:          msg,
				BodyFormatted: msg,
				SendAt:        time.UnixMilli(testEvent().Event.Timestamp),
				Incoming:      true,
				Type:          matrixdb.MessageTypeNewEvent,
			}).Return(nil, nil)

			for _, reaction := range message.ReminderRequestReactions {
				msngr.EXPECT().SendReactionAsync(&messenger.Reaction{
					Reaction:                  reaction,
					MessageExternalIdentifier: testEvent().Event.ID.String(),
					ChannelExternalIdentifier: testEvent().Room.RoomID,
				})
			}

			msngr.EXPECT().SendResponse(gomock.Any()).Return(nil, nil)

			// TODO EXPECT NewMessage

			// Execute
			action.HandleEvent(testEvent(
				withBody(
					msg,
					msg,
				),
			))
		})
	}
	time.Sleep(time.Millisecond * 50)
}

type eventMatcher struct {
	evt *database.Event
}

func (matcher *eventMatcher) Matches(x interface{}) bool {
	evt, ok := x.(*database.Event)
	if !ok {
		return false
	}

	if matcher.evt.ID != 0 {
		if matcher.evt.ID != evt.ID ||
			matcher.evt.CreatedAt != evt.CreatedAt ||
			matcher.evt.UpdatedAt != evt.UpdatedAt {
			return false
		}
	}

	if evt.Time.IsZero() {
		return false
	}

	if matcher.evt.Duration != evt.Duration ||
		matcher.evt.Message != evt.Message ||
		matcher.evt.Active != evt.Active ||
		matcher.evt.ChannelID != evt.ChannelID {
		return false
	}

	if matcher.evt.RepeatInterval != nil {
		if *matcher.evt.RepeatInterval != *evt.RepeatInterval {
			return false
		}
	} else if evt.RepeatInterval != nil {
		return false
	}
	if matcher.evt.RepeatUntil != nil {
		if *matcher.evt.RepeatUntil != *evt.RepeatUntil {
			return false
		}
	} else if evt.RepeatUntil != nil {
		return false
	}
	if matcher.evt.InputID != nil {
		if *matcher.evt.InputID != *evt.InputID {
			return false
		}
	} else if evt.InputID != nil {
		return false
	}

	return true
}

func (matcher *eventMatcher) String() string {
	return fmt.Sprint(matcher.evt)
}
