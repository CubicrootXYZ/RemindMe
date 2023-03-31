package database

import "gorm.io/gorm"

func (service *service) NewIcalInput(input *IcalInput) (*IcalInput, error) {
	err := service.db.Save(input).Error

	return input, err
}

func (service *service) DeleteIcalInput(id uint) error {
	result := service.db.Delete(&IcalInput{Model: gorm.Model{ID: id}})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
