package tests

import (
	"fmt"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
)

type EventMatcher struct {
	evt *database.Event
}

func NewEventMatcher(evt *database.Event) *EventMatcher {
	return &EventMatcher{
		evt: evt,
	}
}

func (matcher *EventMatcher) Matches(x interface{}) bool {
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

func (matcher *EventMatcher) String() string {
	return fmt.Sprint(matcher.evt)
}
