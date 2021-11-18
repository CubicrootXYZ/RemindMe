package database

import "gorm.io/gorm"

// Blocklist holds data about blocked users
type Blocklist struct {
	gorm.Model
	UserIdentifier string `gorm:"uniqueIndex,type:varchar(500)"`
	Reason         string `gorm:"type:text"`
}

// GET DATA

// IsUserBlocked checks if a user is on the blocklist
func (d *Database) IsUserBlocked(userID string) (bool, error) {
	var blocklist Blocklist
	err := d.db.First(&blocklist, "user_identifier = ?", userID).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, err
}

// INSERT DATA

// AddUserToBlocklist blocks the given user with the given reason. The reason is only for internal usecases.
func (d *Database) AddUserToBlocklist(userID string, reason string) error {
	if isBlocked, err := d.IsUserBlocked(userID); isBlocked && err == nil {
		return nil
	}

	blocklist := &Blocklist{
		UserIdentifier: userID,
		Reason:         reason,
	}

	return d.db.Save(blocklist).Error
}
