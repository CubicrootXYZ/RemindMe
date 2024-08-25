package messenger

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/format"
)

const (
	defaultSyncMessageRetries     = 3
	defaultSyncMessageRetryDelay  = time.Second * 5
	defaultAsyncMessageRetries    = 10
	defaultAsyncMessageRetryDelay = time.Second * 10
)

type Response struct {
	Message                   string
	MessageFormatted          string
	RespondToMessage          string
	RespondToMessageFormatted string
	RespondToUserID           string
	RespondToEventID          string
	ChannelExternalIdentifier string
}

func (response *Response) GetResponseMessage() (message, messageFormatted string) {
	message = fmt.Sprintf(
		"> <%s>%s\n\n%s",
		response.RespondToUserID,
		format.StripReply(response.RespondToMessage),
		response.Message,
	)

	quotedMessage := format.StripReplyFormatted(response.RespondToMessageFormatted)
	if response.RespondToMessageFormatted == "" {
		quotedMessage = format.StripReply(response.RespondToMessage)
	}

	messageFormatted = fmt.Sprintf(
		"<mx-reply><blockquote><a href=\"https://matrix.to/#/%s/%s?via=%s\">In reply to</a> <a href=\"https://matrix.to/#/%s\">%s</a><br>%s</blockquote></mx-reply>%s",
		response.ChannelExternalIdentifier,
		response.RespondToEventID,
		format.GetHomeserverFromUserID(response.RespondToUserID),
		response.RespondToUserID,
		response.RespondToUserID,
		quotedMessage,
		response.MessageFormatted,
	)

	return message, messageFormatted
}

func (response *Response) toEvent() *messageEvent {
	message, messageFormatted := response.GetResponseMessage()
	matrixMessage := &messageEvent{
		Body:          message,
		FormattedBody: messageFormatted,
		MsgType:       messageTypeText,
		Format:        formatCustomHTML,
	}
	matrixMessage.RelatesTo.InReplyTo = &struct {
		EventID string `json:"event_id,omitempty"`
	}{EventID: response.RespondToEventID}

	return matrixMessage
}

func PlainTextResponse(msg, replyToEventID, replyToMessage, replyToUser, channelIdentifier string) *Response {
	return &Response{
		Message:                   msg,
		MessageFormatted:          msg,
		RespondToMessage:          replyToMessage,
		RespondToMessageFormatted: replyToMessage,
		RespondToEventID:          replyToEventID,
		RespondToUserID:           replyToUser,
		ChannelExternalIdentifier: channelIdentifier,
	}
}

// SendResponseAsync sends the given response via matrix without blocking the current thread.
// If you need the MessageResponse use SendMessage.
func (messenger *service) SendResponseAsync(response *Response) error {
	go func() {
		_, _ = messenger.sendMessage(response.toEvent(), response.ChannelExternalIdentifier, defaultAsyncMessageRetries, time.Second*10)
	}()

	return nil
}

// SendResponse sends the given response via matrix.
// This will wait for rate limits to expire, thus the request can take some time.
func (messenger *service) SendResponse(response *Response) (*MessageResponse, error) {
	return messenger.sendMessage(response.toEvent(), response.ChannelExternalIdentifier, defaultSyncMessageRetries, time.Second*5)
}
