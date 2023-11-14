package usecase

import (
	dto "github.com/savinnsk/prototype_bot_whatsapp/internal/domain/dto"
	"github.com/savinnsk/prototype_bot_whatsapp/internal/entity"
	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
)

func ProcessNewSchedule(data dto.SaveNewUserAndSchedule) string {

	userFound, _ := gorm.FindUserByPhone(data.CreateUserDto.Phone)
	var userResult *entity.User
	errValidator := ValidateTimeAndDate(data.ScheduleTime, data.ScheduleDate)
	if errValidator != "ok" {
		return errValidator
	}

	if userFound == nil {
		user, err := CreateUserAndReturn(data.CreateUserDto)
		if err != nil {
			print("Error at CreateUserAndReturn Line 11")
			return "🤔 Ops! algo errado ao cadastrar seu nome"
		}
		userResult = user
	} else {
		userResult = userFound

	}

	schedule, err := gorm.FindScheduleByTime(data.ScheduleTime)
	if err != nil {
		print("Error at FindScheduleByTime Line 17")
		return "🤔 Ops! algo errado ao selecionar Horário, verifique se ele esta realmente disponível"
	}

	userSchedule := dto.CreateUserSchedule{
		UserId:     userResult.Id,
		ScheduleId: schedule.Id,
		Time:       schedule.Time,
		Date:       data.ScheduleDate,
	}

	CreateNewSchedule(userSchedule)

	return "ok"

}
