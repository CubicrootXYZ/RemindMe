package formater

import "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"

// EqMessageType tests for equality of two message type lists
func EqMessageType(a, b []database.MessageType) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
