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

func (service *service) RemoveDanglingUsers() (int64, error) {
	res := service.db.Exec(`
DELETE 
FROM matrix_users 
WHERE id IN (
	SELECT DISTINCT(sub.id)
	FROM (SELECT id FROM matrix_users) as sub
	LEFT JOIN matrix_rooms_matrix_users as mrmu ON mrmu.matrix_user_id = sub.id
	WHERE mrmu.matrix_room_id IS NULL
)`)

	return res.RowsAffected, res.Error
}
