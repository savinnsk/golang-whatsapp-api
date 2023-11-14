package entity

type Schedule struct {
	Id        string `gorm:"primaryKey"`
	Time      string `gorm:"not null"`
	DaysWeek  string `gorm:"null"`
	Details   string `gorm:"null"`
	Weekend   bool   `gorm:"null"`
	Holiday   bool   `gorm:"null"`
	Available bool   `gorm:"not null"`
	Disabled  bool   `gorm:"not null"`
}

func NewSchedule(time string, date string, available bool, disabled bool, id string) *Schedule {
	return &Schedule{
		Id:        id,
		Time:      time,
		Available: available,
		Disabled:  disabled,
	}
}
