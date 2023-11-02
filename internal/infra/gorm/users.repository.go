package infra

import entity "github.com/savinnsk/prototype_bot_whatsapp/internal/entity"

func CreateUser(user *entity.User) {

	connection := Init()
	err := connection.db.Create(user).Error
	if err != nil {
		println("Error to create user")
	}
}

func FindUserByPhone(phone string) (*entity.User, error) {
	connection := Init()
	var user entity.User
	err := connection.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUser(user *entity.User) {
	connection := Init()
	connection.db.Delete(&user)
}

func UpdateUser(user *entity.User) {
	connection := Init()
	connection.db.Save(&user)
}

func FindAllUsers() []entity.User {
	connection := Init()
	var users []entity.User
	connection.db.Find(&users)
	return users
}

func LoadUserSchedules(user *entity.User) ([]entity.Schedule, error) {
	connection := Init()
	var schedules []entity.Schedule
	if err := connection.db.Model(user).Preload("Schedules").Find(user).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}
