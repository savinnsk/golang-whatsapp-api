package entity

type Schedule struct {
	Id        int    `gorm:"primaryKey"`
	Time      string `gorm:"not null"`
	Date      *string
	Available bool `gorm:"not null"`
	Disabled  bool `gorm:"not null"`
}

func NewSchedule(time string, date *string, available bool, disabled bool) *Schedule {
	return &Schedule{
		Time:      time,
		Date:      date,
		Available: available,
		Disabled:  disabled,
	}
}
