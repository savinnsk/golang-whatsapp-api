package usecase

import (
	dto "github.com/savinnsk/prototype_bot_whatsapp/internal/domain/dto"
	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
)

func ProcessNewSchedule(data dto.SaveNewUserAndSchedule) string {

	errValidator := ValidateTimeAndDate(data.ScheduleTime, data.ScheduleDate)
	if errValidator != "ok" {
		return errValidator
	}

	user, err := CreateUserAndReturn(data.CreateUserDto)
	if err != nil {
		print("Error at CreateUserAndReturn Line 11")
		return "🤔 Ops! algo errado ao cadastrar seu nome"
	}

	schedule, err := gorm.FindScheduleByTime(data.ScheduleTime)
	if err != nil {
		print("Error at FindScheduleByTime Line 17")
		return "🤔 Ops! algo errado ao selecionar Horário, verifique se ele esta realmente disponível"
	}

	userSchedule := dto.CreateUserSchedule{
		UserId:     user.Id,
		ScheduleId: schedule.Id,
		Time:       schedule.Time,
		Date:       data.ScheduleDate,
	}

	CreateNewSchedule(userSchedule)

	return "ok"

}
