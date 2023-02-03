package database

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("not found")
)

// Service offers an interface for a matrix related database.
type Service interface {
	GetRoomByID(roomID string) (*MatrixRoom, error)
	GetRoomCount() (int64, error)
	NewRoom(room *MatrixRoom) (*MatrixRoom, error)
	UpdateRoom(room *MatrixRoom) (*MatrixRoom, error)
	DeleteRoom(roomID uint) error

	GetUserByID(userID string) (*MatrixUser, error)
	NewUser(user *MatrixUser) (*MatrixUser, error)

	GetLastMessage() (*MatrixMessage, error)
	GetMessageByID(messageID string) (*MatrixMessage, error)
	NewMessage(message *MatrixMessage) (*MatrixMessage, error)
	DeleteAllMessagesFromRoom(roomID uint) error

	GetEventByID(eventID string) (*MatrixEvent, error)
	NewEvent(event *MatrixEvent) (*MatrixEvent, error)
	DeleteAllEventsFromRoom(roomID uint) error
}

// MatrixRoom holds information about a room.
type MatrixRoom struct {
	gorm.Model                   // numeric ID required to match main database in- and outputs
	RoomID          string       `gorm:"unique"`
	Users           []MatrixUser `gorm:"many2many:matrix_rooms_matrix_users;"`
	LastCryptoEvent string
	// TODO somehow get roles back
}

// MatrixUser holds information about an user.
type MatrixUser struct {
	ID      string       `gorm:"primary,size:255"`
	Rooms   []MatrixRoom `gorm:"many2many:matrix_rooms_matrix_users;"`
	Blocked bool
}

type MatrixMessageType string

var (
	MessageTypeWelcome = MatrixMessageType("WELCOME")
)

// MatrixMessage holds information about a matrix message.
type MatrixMessage struct {
	ID            string `gorm:"primary"`
	UserID        string `gorm:"size:255"`
	User          MatrixUser
	RoomID        uint
	Room          MatrixRoom
	Body          string
	BodyFormatted string
	SendAt        time.Time
	Type          MatrixMessageType
	Incoming      bool
}

// MatrixEvent holds information about a state event (e.g. leave, join).
type MatrixEvent struct {
	ID     string `gorm:"primary,size:255"`
	UserID string `gorm:"size:255"`
	User   MatrixUser
	RoomID uint
	Room   MatrixRoom
	Type   string
	SendAt time.Time
}