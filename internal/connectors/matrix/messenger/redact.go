package messenger

import (
	"errors"
	"fmt"
	"time"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

type Redact struct {
	RoomID    string
	MessageID string
}

func (messenger *service) SendRedactAsync(redact *Redact) error {
	go func() {
		_ = messenger.sendRedact(redact.RoomID, redact.MessageID, 10, time.Second*15)
	}()

	return nil
}

func (messenger *service) sendRedact(roomID string, eventID string, retries uint, retryTime time.Duration) error {
	var err error
	maxRetries := retries

	for retries > 0 {
		// Wait until the rate limit is gone again
		for time.Until(messenger.state.rateLimitedUntil) >= 0 {
			time.Sleep(time.Second * 5)
			continue
		}

		_, err := messenger.client.RedactEvent(id.RoomID(roomID), id.EventID(eventID))
		if err == nil {
			// No error, fine return the result
			return nil
		} else if errors.Is(err, mautrix.MLimitExceeded) {
			// Rate limit is exceeded so wait until we can send requests again
			messenger.encounteredRateLimit()
			messenger.logger.Info("sending message is stopped since we ran in a rate limit")
			continue
		} else if errors.Is(err, mautrix.MForbidden) || errors.Is(err, mautrix.MUnknownToken) || errors.Is(err, mautrix.MMissingToken) || errors.Is(err, mautrix.MBadJSON) || errors.Is(err, mautrix.MNotJSON) || errors.Is(err, mautrix.MUnsupportedRoomVersion) || errors.Is(err, mautrix.MIncompatibleRoomVersion) {
			// Errors indicating that the request is invalid, do not try again
			messenger.logger.Info("sending message failed", "error", err)
			return err
		}

		messenger.logger.Info(fmt.Sprintf("Sending message failed in try %d from try %d with error: %s", retries, maxRetries, err.Error()))

		retries--
		time.Sleep(retryTime * (time.Duration(maxRetries) - time.Duration(retries)))
	}

	if err == nil {
		err = ErrRetriesExceeded
	}

	messenger.logger.Info("sending message failed and retries are exceeded", "error", err)
	return err
}
