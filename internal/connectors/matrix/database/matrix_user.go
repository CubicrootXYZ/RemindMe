package database

import (
	"errors"

	"gorm.io/gorm"
)

func (service *service) GetUserByID(userID string) (*MatrixUser, error) {
	var user MatrixUser
	err := service.db.Preload("Rooms").First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, err
}

func (service *service) NewUser(user *MatrixUser) (*MatrixUser, error) {
	err := service.db.Create(user).Error

	return user, err
}
