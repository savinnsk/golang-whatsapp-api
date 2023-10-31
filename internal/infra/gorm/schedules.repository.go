package infra

import entity "github.com/savinnsk/prototype_bot_whatsapp/internal/entity"

func CreateSchedule(schedule *entity.Schedule) error {
	connection := Init()
	return connection.db.Create(schedule).Error
}

func UpdateSchedule(schedule *entity.Schedule) error {
	connection := Init()
	return connection.db.Save(schedule).Error
}

func DeleteSchedule(schedule *entity.Schedule) error {
	connection := Init()
	return connection.db.Delete(schedule).Error
}

func LoadAllSchedules() ([]entity.Schedule, error) {
	connection := Init()
	var schedules []entity.Schedule
	if err := connection.db.Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}
