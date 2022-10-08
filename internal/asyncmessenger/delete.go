package asyncmessenger

import (
	"errors"
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

// Delete holds information about a delete request
type Delete struct {
	ExternalIdentifier        string
	ChannelExternalIdentifier string
}

// DeleteMessageAsync removes a message from the channel
func (messenger *messenger) DeleteMessageAsync(delete *Delete) error {
	go messenger.deleteMessage(delete, 10, 15*time.Second)
	return nil
}

func (messenger *messenger) deleteMessage(delete *Delete, retries uint, retryTime time.Duration) error {
	var err error
	maxRetries := retries

	for retries > 0 {
		// Wait until the rate limit is gone again
		for time.Until(messenger.state.rateLimitedUntil) >= 0 {
			time.Sleep(time.Second * 5)
			continue
		}

		_, err := messenger.client.RedactEvent(id.RoomID(delete.ChannelExternalIdentifier), id.EventID(delete.ExternalIdentifier))
		if err == nil {
			// No error, fine return the result
			return nil
		} else if errors.Is(err, mautrix.MLimitExceeded) {
			// Rate limit is exceeded so wait until we can send requests again
			messenger.encounteredRateLimit()
			log.Info("Deleting message is stopped since we ran in a rate limit")
			continue
		} else if errors.Is(err, mautrix.MForbidden) || errors.Is(err, mautrix.MUnknownToken) || errors.Is(err, mautrix.MMissingToken) || errors.Is(err, mautrix.MBadJSON) || errors.Is(err, mautrix.MNotJSON) || errors.Is(err, mautrix.MUnsupportedRoomVersion) || errors.Is(err, mautrix.MIncompatibleRoomVersion) {
			// Errors indicating that the request is invalid, do not try again
			log.Info("Deleting message failed with error: " + err.Error())
			return err
		} else {
			log.Info(fmt.Sprintf("Deleting message failed in try %d from try %d with error: %s", retries, maxRetries, err.Error()))
		}

		retries--
		time.Sleep(retryTime * (time.Duration(maxRetries) - time.Duration(retries)))
	}

	if err == nil {
		err = ErrRetriesExceeded
	}

	log.Info("Deleting message failed and retries are exceeded. Error is: " + err.Error())
	return err
}
