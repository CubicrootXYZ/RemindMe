package formater

import (
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestFormater_EqMessageType(t *testing.T) {
	lists := make([][]database.MessageType, 0)
	lists = append(lists, []database.MessageType{})
	lists = append(lists, nil)
	lists = append(lists, []database.MessageType{database.MessageTypeDailyReminderDeleteSuccess})
	lists = append(lists, []database.MessageType{database.MessageTypeDailyReminderDeleteSuccess, database.MessageTypeIcalRenew})
	lists = append(lists, []database.MessageType{database.MessageTypeIcalRenew, database.MessageTypeDailyReminderDeleteSuccess})
	lists = append(lists, []database.MessageType{database.MessageTypeDailyReminderDeleteSuccess, database.MessageTypeIcalRenew, database.MessageTypeIcalLink, database.MessageTypeReminderFail, database.MessageTypeReminderUpdateSuccess})

	for i, list := range lists {
		for j, compList := range lists {
			match := EqMessageType(list, compList)

			if i == j {
				assert.Truef(t, match, "List %d and %d should match", i, j)
			} else {
				assert.False(t, match, "List %d and %d should not match", i, j)
			}
		}
	}
}
