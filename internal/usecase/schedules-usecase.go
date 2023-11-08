package usecase

import (
	"strconv"

	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
)

func FilterSchedules() []string {
	// pending logic to deal with dates
	var scheduleArray []string
	schedules, _ := gorm.LoadAllSchedules()

	for _, schedule := range schedules {
		result := ValidateTimeIsCurrent(schedule.Time)
		if result == "ok" {
			scheduleArray = append(scheduleArray, schedule.Time)
		}

	}

	return scheduleArray
}

func VerifyScheduleBasedAtArray(choice string, schedulesFiltered []string) string {

	choiceInt, err := strconv.Atoi(choice)
	if err != nil {
		return "Invalid choice: not a number"
	}
	newChose := choiceInt - 2
	return schedulesFiltered[newChose]
}
