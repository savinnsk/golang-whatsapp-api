package entity

type UserSchedule struct {
	Id         int    `gorm:"primaryKey"`
	UserID     string `gorm:"not nul"`
	ScheduleID string `gorm:"not null"`
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
