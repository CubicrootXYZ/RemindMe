package messenger

import (
	"errors"
	"time"

	"maunium.net/go/mautrix"
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
		{
			Body:     message.Body,
			Mimetype: mimetypeTextPlain,
		},
		{
			Body:     message.BodyHTML,
			Mimetype: mimetypeTextHTML,
		},
	}

	if message.ResponseToMessage != "" {
		messageEvent.RelatesTo.InReplyTo = &struct {
			EventID string `json:"event_id,omitempty"`
		}{EventID: message.ResponseToMessage}
	}

	return &messageEvent
}

// PlainTextMessage constructs a simple plain text message
func PlainTextMessage(message string, channelIdentifier string) *Message {
	return &Message{
		Body:                      message,
		BodyHTML:                  message,
		ChannelExternalIdentifier: channelIdentifier,
	}
}

// HTMLMessage constructs a simple HTML message with a plaintext fallback
func HTMLMessage(plaintextMessage, htmlMessage, channelIdentifier string) *Message {
	return &Message{
		Body:                      plaintextMessage,
		BodyHTML:                  htmlMessage,
		ChannelExternalIdentifier: channelIdentifier,
	}
}

type MessageResponse struct {
	ExternalIdentifier string
	Timestamp          time.Time
}

// SendMessageAsync sends the given message via matrix without blocking the current thread.
// If you need the MessageResponse use SendMessage.
func (messenger *service) SendMessageAsync(message *Message) error {
	go func() {
		_, _ = messenger.sendMessage(
			message.toEvent(),
			message.ChannelExternalIdentifier,
			defaultAsyncMessageRetries,
			defaultAsyncMessageRetryDelay,
		)
	}()

	return nil
}

// SendMessage sends the given message via matrix.
// This will wait for rate limits to expire, thus the request can take some time.
func (messenger *service) SendMessage(message *Message) (*MessageResponse, error) {
	return messenger.sendMessage(
		message.toEvent(),
		message.ChannelExternalIdentifier,
		defaultSyncMessageRetries,
		defaultSyncMessageRetryDelay,
	)
}

// sendMessage will take care of sending the message via matrix
// The message sending will be tried for retries times and the time between retries is retry * retryTime
func (messenger *service) sendMessage(messageEvent *messageEvent, channel string, retries uint, retryTime time.Duration) (*MessageResponse, error) {
	var err error
	maxRetries := retries

	messenger.metricEventOutCount.
		WithLabelValues("message").
		Inc()

	for retries > 0 {
		// Wait until the rate limit is gone again
		for time.Until(messenger.state.rateLimitedUntil) >= 0 {
			time.Sleep(time.Second * 5)
			continue
		}

		response, err := messenger.sendMessageEvent(messageEvent, channel, messageEvent.getEventType())
		if err == nil {
			// No error, fine return the result
			return messenger.mautrixRespSendEventToMessageResponse(response), nil
		} else if errors.Is(err, mautrix.MLimitExceeded) {
			// Rate limit is exceeded so wait until we can send requests again
			messenger.encounteredRateLimit()
			messenger.logger.Info("sending message is stopped since we ran in a rate limit")
			continue
		} else if errors.Is(err, mautrix.MForbidden) || errors.Is(err, mautrix.MUnknownToken) || errors.Is(err, mautrix.MMissingToken) || errors.Is(err, mautrix.MBadJSON) || errors.Is(err, mautrix.MNotJSON) || errors.Is(err, mautrix.MUnsupportedRoomVersion) || errors.Is(err, mautrix.MIncompatibleRoomVersion) {
			// Errors indicating that the request is invalid, do not try again
			messenger.logger.Info("sending message failed", "error", err)
			return nil, err
		}

		messenger.logger.Info("sending message failed", "try", retries, "max_tries", maxRetries, "error", err)

		retries--
		time.Sleep(retryTime * (time.Duration(maxRetries) - time.Duration(retries)))
	}

	if err == nil {
		err = ErrRetriesExceeded
	}

	messenger.logger.Info("sending message failed and retries are exceeded", "error", err)
	return nil, err
}

func (messenger *service) mautrixRespSendEventToMessageResponse(responseEvent *mautrix.RespSendEvent) *MessageResponse {
	return &MessageResponse{
		ExternalIdentifier: responseEvent.EventID.String(),
		Timestamp:          time.Now(), // Unfortunately the response event does not contain a timestamp
	}
}
