package usecase

import (
	dto "github.com/savinnsk/prototype_bot_whatsapp/internal/domain/dto"
	"github.com/savinnsk/prototype_bot_whatsapp/internal/entity"
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

func LoadAllUserSchedules(phone string) ([]entity.UserSchedule, error) {
	user, err := gorm.FindUserByPhone(phone)

	if err != nil {
		return nil, err
	}
	schedules, err := gorm.LoadUserSchedulesByUserID(user.Id)

	if err != nil {
		println("Error LoadAllUserSchedules Schedule Line 24")
		return nil, err
	}

	return schedules, nil
}
