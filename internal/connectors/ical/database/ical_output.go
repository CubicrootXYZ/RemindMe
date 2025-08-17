package database

import (
	"errors"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/random"
	"gorm.io/gorm"
)

func (service *service) NewIcalOutput(output *IcalOutput) (*IcalOutput, error) {
	output.Token = random.URLSaveString(random.Intn(10) + 30)

	err := service.db.Save(output).Error

	return output, err
}

func (service *service) GetIcalOutputByID(id uint) (*IcalOutput, error) {
	var entity IcalOutput

	err := service.db.First(&entity, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}

	return &entity, err
}

func (service *service) GenerateNewToken(output *IcalOutput) (*IcalOutput, error) {
	// New token with 30-40 characters length.
	output.Token = random.URLSaveString(random.Intn(10) + 30)

	err := service.db.Save(output).Error

	return output, err
}

func (service *service) DeleteIcalOutput(id uint) error {
	result := service.db.Delete(&IcalOutput{Model: gorm.Model{ID: id}})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
