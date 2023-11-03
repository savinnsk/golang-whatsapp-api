package usecase

import (
	dto "github.com/savinnsk/prototype_bot_whatsapp/internal/domain/dto"
	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
)

func CreateNewSchedule(data dto.CreateUserSchedule) error {
	err := gorm.CreateUserSchedule(data)

	if err != nil {
		println("Error CreateNewSchedule Schedule Line 12")
		return err
	}

	return nil
}
