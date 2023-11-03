package usecase

import (
	entity "github.com/savinnsk/prototype_bot_whatsapp/internal/entity"
	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
)

func LoadAllValidSchedulesDates() []entity.Schedule {
	// pending logic to deal with dates
	schedules, err := gorm.LoadAllSchedules()

	if err != nil {
		return nil
	}

	return schedules
}
