package asyncmessenger

import (
	"errors"
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

// Message holds information about a message
type Message struct {
	Body                      string // Plain text message use \n for newlines
	BodyHTML                  string // HTML formatted message - optional
	ResponseToMessage         string // ID of a message to respond to - optional
	ChannelExternalIdentifier string // Channel ID to send the message to
}

// toEvent converts a Message struct into a matrix message event
func (message *Message) toEvent() *messageEvent {
	messageEvent := messageEvent{
		Body:          message.Body,
		FormattedBody: message.BodyHTML,
		Format:        formatCustomHTML,
		MsgType:       messageTypeText,
		Type:          eventTypeRoomMessage,
	}
	messageEvent.MSC1767Message = []matrixMSC1767Event{
		matrixMSC1767Event{
			Body:     message.Body,
			Mimetype: mimetypeTextPlain,
		},
		matrixMSC1767Event{
			Body:     message.BodyHTML,
			Mimetype: mimetypeTextHTML,
		},
	}

	if message.ResponseToMessage != "" {
		messageEvent.RelatesTo.InReplyTo.EventID = message.ResponseToMessage
	}

	return &messageEvent
}

type MessageResponse struct {
	ExternalIdentifier string
	Timestamp          int64
}

// SendMessageAsync sends the given message via matrix without blocking the current thread.
// If you need the MessageResponse use SendMessage.
func (messenger *messenger) SendMessageAsync(message *Message) error {
	go messenger.sendMessage(message.toEvent(), message.ChannelExternalIdentifier, 10, time.Second*10)

	return nil
}

// SendMessage sends the given message via matrix.
// This will wait for rate limits to expire, thus the request can take some time.
func (messenger *messenger) SendMessage(message *Message) (*MessageResponse, error) {
	return messenger.sendMessage(message.toEvent(), message.ChannelExternalIdentifier, 3, time.Second*5)
}

// sendMessage will take care of sending the message via matrix
// The message sending will be tried for retries times and the time between retries is retry * retryTime
func (messenger *messenger) sendMessage(messageEvent *messageEvent, channel string, retries uint, retryTime time.Duration) (*MessageResponse, error) {
	var err error
	maxRetries := retries

	for retries > 0 {
		// Wait until the rate limit is gone again
		for time.Until(messenger.state.rateLimitedUntil) >= 0 {
			time.Sleep(time.Second * 5)
			continue
		}

		response, err := messenger.sendMessageEvent(messageEvent, channel, event.EventMessage)
		if err == nil {
			// No error, fine return the result
			return messenger.mautrixRespSendEventToMessageResponse(response), nil
		} else if errors.Is(err, mautrix.MLimitExceeded) {
			// Rate limit is exceeded so wait until we can send requests again
			messenger.encounteredRateLimit()
			log.Info("Sending message is stopped since we ran in a rate limit")
			continue
		} else if errors.Is(err, mautrix.MForbidden) || errors.Is(err, mautrix.MUnknownToken) || errors.Is(err, mautrix.MMissingToken) || errors.Is(err, mautrix.MBadJSON) || errors.Is(err, mautrix.MNotJSON) || errors.Is(err, mautrix.MUnsupportedRoomVersion) || errors.Is(err, mautrix.MIncompatibleRoomVersion) {
			// Errors indicating that the request is invalid, do not try again
			log.Info("Sending message failed with error: " + err.Error())
			return nil, err
		} else {
			log.Info(fmt.Sprintf("Sending message failed in try %d from try %d with error: %s", retries, maxRetries, err.Error()))
		}

		retries--
		time.Sleep(retryTime * (time.Duration(maxRetries) - time.Duration(retries)))
	}

	if err == nil {
		err = ErrRetriesExceeded
	}

	log.Info("Sending message failed and retries are exceeded. Error is: " + err.Error())
	return nil, err
}

func (messenger *messenger) mautrixRespSendEventToMessageResponse(responseEvent *mautrix.RespSendEvent) *MessageResponse {
	return &MessageResponse{
		ExternalIdentifier: responseEvent.EventID.String(),
		Timestamp:          time.Now().Unix(), // Unfortunately the response event does not contain a timestamp
	}
}
