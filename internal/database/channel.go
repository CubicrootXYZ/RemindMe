package database

import (
	"errors"

	"gorm.io/gorm"
)

func (service *service) NewChannel(channel *Channel) (*Channel, error) {
	err := service.db.Create(&channel).Error
	return channel, err
}

func (service *service) GetChannelByID(channelID uint) (*Channel, error) {
	var channel Channel
	err := service.db.Preload("Inputs").Preload("Outputs").First(&channel, "channels.id = ?", channelID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &channel, err
}

func (service *service) GetChannels() ([]Channel, error) {
	var channels []Channel
	err := service.db.Find(&channels).Error

	return channels, err
}

func (service *service) AddInputToChannel(channelID uint, input *Input) error {
	input.ChannelID = channelID
	return service.db.Save(input).Error
}

func (service *service) AddOutputToChannel(channelID uint, output *Output) error {
	output.ChannelID = channelID
	return service.db.Save(output).Error
}

func (service *service) RemoveInputFromChannel(channelID, inputID uint) error {
	tx := service.newSession()
	input, err := tx.GetInputByID(inputID)
	if err != nil {
		return tx.rollbackWithError(err)
	}

	// Inform the input service that we will delete the input
	inputService, ok := service.config.InputServices[input.InputType]
	if !ok {
		return tx.rollbackWithError(ErrUnknownInput)
	}

	err = inputService.InputRemoved(input.InputType, input.InputID)
	if err != nil {
		return tx.rollbackWithError(err)
	}

	// Delete the input permanently
	err = tx.deleteInput(channelID, inputID)
	if err != nil {
		return tx.rollbackWithError(err)
	}

	err = tx.commit()
	if err != nil {
		return tx.rollbackWithError(err)
	}

	return nil
}

func (service *service) deleteInput(channelID, inputID uint) error {
	return service.db.Unscoped().Where("id = ? AND channel_id = ?", inputID, channelID).Delete(&Input{}).Error
}

func (service *service) RemoveOutputFromChannel(channelID, outputID uint) error {
	tx := service.newSession()
	output, err := tx.GetOutputByID(outputID)
	if err != nil {
		return tx.rollbackWithError(err)
	}

	// Inform the output service that we will delete the output
	outputService, ok := service.config.OutputServices[output.OutputType]
	if !ok {
		return tx.rollbackWithError(ErrUnknownOutput)
	}

	err = outputService.OutputRemoved(output.OutputType, output.OutputID)
	if err != nil {
		return tx.rollbackWithError(err)
	}

	// Delete the output permanently
	err = tx.deleteOutput(channelID, outputID)
	if err != nil {
		return tx.rollbackWithError(err)
	}

	err = tx.commit()
	if err != nil {
		return tx.rollbackWithError(err)
	}

	return nil
}

func (service *service) deleteOutput(channelID, outputID uint) error {
	return service.db.Unscoped().Where("id = ? AND channel_id = ?", outputID, channelID).Delete(&Output{}).Error
}

func (service *service) UpdateChannel(channel *Channel) (*Channel, error) {
	err := service.db.Save(channel).Error

	return channel, err
}

func (service *service) DeleteChannel(channelID uint) error {
	result := service.db.Unscoped().Delete(&Channel{}, "id = ?", channelID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return ErrNotFound
	}

	return nil
}
