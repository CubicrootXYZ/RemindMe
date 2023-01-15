package database

import (
	"errors"

	"gorm.io/gorm"
)

func (service *service) GetOutputByID(outputID uint) (*Output, error) {
	var output Output
	err := service.db.First(&output, "id = ?", outputID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return &output, err
}
