package messenger

import (
	"time"

	"maunium.net/go/mautrix/id"
)

// roomCache is a short term cache to avoid querying for room members to often
type roomCache map[id.RoomID]roomCacheEntry

type roomCacheEntry struct {
	CachedAt    time.Time
	RoomMembers []id.UserID
}

// GetUsers returns users in a room if stored in the cache, otherwise nil
func (cache roomCache) GetUsers(room id.RoomID) []id.UserID {
	if entry, ok := cache[room]; ok {
		if time.Since(entry.CachedAt) > time.Minute*2 {
			delete(cache, room)
			return nil
		}

		return entry.RoomMembers
	}

	return nil
}

// AddUsers adds a room member list to the cache
func (cache roomCache) AddUsers(room id.RoomID, users []id.UserID) {
	cache[room] = roomCacheEntry{
		CachedAt:    time.Now(),
		RoomMembers: users,
	}
}
