package formater

import "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"

// EqMessageType tests for equality of two message type lists (entries must have same order)
func EqMessageType(a, b []database.MessageType) bool {
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

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
