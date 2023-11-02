package entity

type UserSchedule struct {
	UserID     int    `gorm:"primaryKey"`
	ScheduleID int    `gorm:"primaryKey"`
	Date       string `gorm:"null"`
	Time       string `gorm:"null"`
}

func NewUserSchedule(userID, scheduleID int, data, time string) *UserSchedule {
	return &UserSchedule{
		UserID:     userID,
		ScheduleID: scheduleID,
		Date:       data,
		Time:       time,
	}
}
