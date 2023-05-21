package msghelper

import (
	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix"
	matrixdb "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
)

// Storer is a helper for sending and storing messages async.
type Storer struct {
	db        matrixdb.Service
	messenger messenger.Messenger
	logger    gologger.Logger
}

// NewStorer assembles a new storer.
func NewStorer(db matrixdb.Service, messenger messenger.Messenger, logger gologger.Logger) *Storer {
	return &Storer{
		db:        db,
		messenger: messenger,
		logger:    logger,
	}
}

// SendAndStoreMessage outsources message sending and storing. Call it asynchronous with "go SendAndStoreMessage(...)"
func (storer *Storer) SendAndStoreMessage(message, messageFormatted string, messageType matrixdb.MatrixMessageType, event matrix.MessageEvent) {
	resp, err := storer.messenger.SendMessage(messenger.HTMLMessage(message, messageFormatted, event.Room.RoomID))
	if err != nil {
		storer.logger.Err(err)
		return
	}

	sender := event.Event.Sender.String()
	dbMessage := matrixdb.MatrixMessage{
		ID:            resp.ExternalIdentifier,
		Type:          messageType,
		Incoming:      false,
		SendAt:        resp.Timestamp,
		Body:          message,
		BodyFormatted: messageFormatted,
		UserID:        &sender,
		RoomID:        event.Room.ID,
	}

	_, err = storer.db.NewMessage(&dbMessage)
	if err != nil {
		storer.logger.Err(err)
	}
}

// SendAndStoreResponse outsources response sending and storing. Call it asynchronous with "go SendAndStoreMessage(...)"
func (storer *Storer) SendAndStoreResponse(message string, messageType matrixdb.MatrixMessageType, event matrix.MessageEvent) {
	resp, err := storer.messenger.SendResponse(messenger.PlainTextResponse(
		message,
		event.Event.ID.String(),
		event.Content.Body,
		event.Event.Sender.String(),
		event.Room.RoomID,
	))
	if err != nil {
		storer.logger.Err(err)
		return
	}

	sender := event.Event.Sender.String()
	dbMessage := matrixdb.MatrixMessage{
		ID:            resp.ExternalIdentifier,
		Type:          messageType,
		Incoming:      false,
		SendAt:        resp.Timestamp,
		Body:          message,
		BodyFormatted: message,
		UserID:        &sender,
		RoomID:        event.Room.ID,
	}

	_, err = storer.db.NewMessage(&dbMessage)
	if err != nil {
		storer.logger.Err(err)
	}
}
