package asyncmessenger

import (
	"fmt"
	"time"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/formater"
)

type Response struct {
	Message                   string
	MessageFormatted          string
	RespondToMessage          string
	RespondToMessageFormatted string
	RespondToUserID           string
	RoomID                    string
	RespondToEventID          string
	ChannelExternalIdentifier string
}

func (response *Response) getResponseMessage() (message, messageFormatted string) {
	message = fmt.Sprintf(
		"> <%s>%s\n\n%s",
		response.RespondToUserID,
		response.RespondToMessage,
		response.Message,
	)

	messageFormatted = fmt.Sprintf(
		"<mx-reply><blockquote><a href=\"https://matrix.to/#/%s/%s?via=%s\">In reply to</a> <a href=\"https://matrix.to/#/%s\">%s</a><br>%s</blockquote></mx-reply>%s",
		response.RoomID,
		response.RespondToEventID,
		formater.GetHomeserverFromUserID(response.RespondToUserID),
		response.RespondToUserID,
		response.RespondToUserID,
		response.RespondToMessageFormatted,
		response.MessageFormatted,
	)

	return message, messageFormatted
}

func (response *Response) toEvent() *messageEvent {
	message, messageFormatted := response.getResponseMessage()
	matrixMessage := &messageEvent{
		Body:          message,
		FormattedBody: messageFormatted,
		MsgType:       messageTypeText,
		Format:        formatCustomHTML,
	}
	matrixMessage.RelatesTo.InReplyTo.EventID = response.RespondToEventID

	return matrixMessage
}

// SendResponseAsync sends the given response via matrix without blocking the current thread.
// If you need the MessageResponse use SendMessage.
func (messenger *messenger) SendResponseAsync(response *Response) error {
	go messenger.sendMessage(response.toEvent(), response.ChannelExternalIdentifier, 10, time.Second*10)

	return nil
}

// SendResponse sends the given response via matrix.
// This will wait for rate limits to expire, thus the request can take some time.
func (messenger *messenger) SendResponse(response *Response) (*MessageResponse, error) {
	return messenger.sendMessage(response.toEvent(), response.ChannelExternalIdentifier, 3, time.Second*5)
}
