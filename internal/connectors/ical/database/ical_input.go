package database

import (
	"errors"

	"gorm.io/gorm"
)

func (service *service) NewIcalInput(input *IcalInput) (*IcalInput, error) {
	err := service.db.Save(input).Error

	return input, err
}

func (service *service) GetIcalInputByID(id uint) (*IcalInput, error) {
	var entity IcalInput

	err := service.db.First(&entity, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return &entity, err
}

func (service *service) ListIcalInputs(opts *ListIcalInputsOpts) ([]IcalInput, error) {
	var inputs []IcalInput

	query := service.db

	if opts.Disabled != nil {
		query = query.Where("disabled = ?", *opts.Disabled)
	}

	err := query.Find(&inputs).Error

	return inputs, err
}

func (service *service) UpdateIcalInput(entity *IcalInput) (*IcalInput, error) {
	err := service.db.Save(entity).Error
	return entity, err
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
