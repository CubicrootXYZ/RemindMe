package asyncmessenger

import (
	"fmt"

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
