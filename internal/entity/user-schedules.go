package entity

type UserSchedule struct {
	UserID     string `gorm:"primaryKey"`
	ScheduleID string `gorm:"primaryKey"`
	Date       string `gorm:"null"`
	Time       string `gorm:"null"`
}

func NewUserSchedule(userID string, scheduleID string, data, time string) *UserSchedule {
	return &UserSchedule{
		UserID:     userID,
		ScheduleID: scheduleID,
		Date:       data,
		Time:       time,
	}
}
