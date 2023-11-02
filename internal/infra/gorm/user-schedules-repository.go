package infra

import entity "github.com/savinnsk/prototype_bot_whatsapp/internal/entity"

func CreateUserSchedule(userID int, scheduleID int, date string, time string) error {
	connection := Init()
	userSchedule := entity.UserSchedule{
		UserID:     userID,
		ScheduleID: scheduleID,
		Date:       date,
		Time:       time,
	}
	err := connection.db.Create(&userSchedule).Error
	if err != nil {
		return err
	}
	return nil
}

func LoadUserSchedulesByUserID(userID int) ([]entity.Schedule, error) {
	connection := Init()
	var userSchedules []entity.Schedule
	err := connection.db.Where("user_id = ?", userID).Find(&userSchedules).Error
	if err != nil {
		return nil, err
	}
	return userSchedules, nil
}

func DeleteUserSchedule(userID int, scheduleID int) error {
	connection := Init()
	err := connection.db.Where("user_id = ? AND schedule_id = ?", userID, scheduleID).Delete(entity.UserSchedule{}).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserSchedule(userID int, scheduleID int, newDate string, newTime string) error {
	connection := Init()

	var userSchedule entity.UserSchedule
	err := connection.db.Where("user_id = ? AND schedule_id = ?", userID, scheduleID).First(&userSchedule).Error
	if err != nil {
		return err
	}

	userSchedule.Date = newDate
	userSchedule.Time = newTime

	err = connection.db.Save(&userSchedule).Error
	if err != nil {
		return err
	}

	return nil
}
