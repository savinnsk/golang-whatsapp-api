package usecase

import (
	"time"
)

func ValidateTimeAndDate(timeSchedule string, dateSchedule string) string {

	currentDate := time.Now().Format("02/01/2006")
	currentTime := time.Now().Format("15:04")

	scheduleDate, _ := time.Parse("02/01/2006", dateSchedule)
	currentDateObj, _ := time.Parse("02/01/2006", currentDate)

	if scheduleDate.Before(currentDateObj) {
		return "😅, Verifique se a *Data* é valida."
	}

	if scheduleDate.Equal(currentDateObj) {
		scheduleTime, _ := time.Parse("15:04", timeSchedule)
		currentTimeObj, _ := time.Parse("15:04", currentTime)

		if scheduleTime.Before(currentTimeObj) {
			return "😅, Verifique se a *Hora* é valida."
		}
	}

	return "ok"
}

func ValidateTimeIsCurrent(timeSchedule string) string {
	currentTime := time.Now().Format("15:04")

	scheduleTime, _ := time.Parse("15:04", timeSchedule)
	currentTimeObj, _ := time.Parse("15:04", currentTime)

	if scheduleTime.Before(currentTimeObj) {
		return "😅, Verifique se a *Hora* é valida."
	}

	return "ok"
}
