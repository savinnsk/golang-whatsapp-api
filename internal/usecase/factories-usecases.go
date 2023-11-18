package usecase

import (
	"fmt"

	dto "github.com/savinnsk/prototype_bot_whatsapp/internal/domain/dto"
	"github.com/savinnsk/prototype_bot_whatsapp/internal/entity"
	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
)

func ProcessNewSchedule(data dto.SaveNewUserAndSchedule) string {
	var userResult *entity.User
	userFound, err := gorm.FindUserByPhone(data.CreateUserDto.Phone)
	if err != nil {
		user, err := CreateUserAndReturn(data.CreateUserDto)

		if err != nil {
			print("🤔  🤔  🤔  🤔  🤔 Error at CreateUserAndReturn Line 11")
			return "🤔 Ops! algo errado ao cadastrar seu nome"
		}

		userResult = user
	}

	errValidator := ValidateTimeAndDate(data.ScheduleTime, data.ScheduleDate)

	fmt.Println(" 🤔  🤔  🤔  🤔  🤔  HAHHAAHHAHHAH>>>>>>", errValidator)
	if errValidator != "ok" {
		return errValidator
	}
	if userFound != nil {

		userResult = userFound
	}

	fmt.Println(" 🤔  🤔  🤔  🤔  🤔  HAHHAAHHAHHAH>>>>>>", userResult)

	schedule, err := gorm.FindScheduleByTime(data.ScheduleTime)
	if err != nil {
		print("🤔  🤔  🤔  🤔  🤔  Error at FindScheduleByTime Line 17")
		return "🤔 Ops! algo errado ao selecionar Horário, verifique se ele esta realmente disponível"
	}
	fmt.Println("🤔  🤔  🤔  🤔  🤔  HAHHAAHHAHHAH>>>>>>", schedule)

	userSchedule := dto.CreateUserSchedule{
		UserId:     userResult.Id,
		ScheduleId: schedule.Id,
		Time:       schedule.Time,
		Date:       data.ScheduleDate,
	}

	err = CreateNewSchedule(userSchedule)
	println("💀: Error factories.usecases.go line 52 : ", err)
	if err != nil {
		println("💀: Error factories.usecases.go line 52 : ", err)
	}

	return "ok"

}
