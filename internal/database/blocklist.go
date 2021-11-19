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

// GetBlockedUserList lists all blocked users
func (d *Database) GetBlockedUserList() ([]Blocklist, error) {
	blocklists := make([]Blocklist, 0)
	err := d.db.Find(&blocklists).Error
	return blocklists, err
}

// INSERT DATA

// AddUserToBlocklist blocks the given user with the given reason. The reason is only for internal usecases.
func (d *Database) AddUserToBlocklist(userID string, reason string) error {
	isBlocked, err := d.IsUserBlocked(userID)
	if isBlocked && err == nil {
		return nil
	} else if err != nil {
		return err
	}

	blocklist := &Blocklist{
		UserIdentifier: userID,
		Reason:         reason,
	}

	return d.db.Save(blocklist).Error
}

// DELETE DATA

// RemoveUserFromBlocklist removes the given user from the blocklist
func (d *Database) RemoveUserFromBlocklist(userID string) error {
	return d.db.Exec("DELETE FROM `blocklists` WHERE user_identifier = ?", userID).Error
}
