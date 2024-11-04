package messenger

import (
	"log/slog"
	"sync"
	"time"

	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type service struct {
	roomUserCache roomCache
	config        *Config
	client        MatrixClient
	db            database.Service
	logger        gologger.Logger
	state         *state
}

type Config struct {
	DisableReactions bool
}

type state struct {
	rateLimitedUntil      time.Time // If we run into a rate limit this will tell us to stop operation
	rateLimitedUntilMutex sync.Mutex
}

func NewMessenger(config *Config, db database.Service, matrixClient MatrixClient, logger *slog.Logger) (Messenger, error) {
	return &service{
		roomUserCache: make(roomCache),
		config:        config,
		client:        matrixClient,
		db:            db,
		// TODO logger:        logger,
		state: &state{
			rateLimitedUntilMutex: sync.Mutex{},
		},
	}, nil
}

// sendMessageEvent sends a message event to matrix, will take care of encryption if available
func (messenger *service) sendMessageEvent(messageEvent *messageEvent, roomID string, eventType event.Type) (*mautrix.RespSendEvent, error) {
	messenger.logger.Infof("Sending message to room %s", roomID)
	return messenger.client.SendMessageEvent(id.RoomID(roomID), eventType, &messageEvent)
}

func (messenger *service) encounteredRateLimit() {
	messenger.state.rateLimitedUntilMutex.Lock()
	messenger.state.rateLimitedUntil = time.Now().Add(time.Minute)
	messenger.state.rateLimitedUntilMutex.Unlock()
}
