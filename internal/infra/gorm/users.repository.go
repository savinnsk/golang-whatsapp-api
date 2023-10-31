package infra

import entity "github.com/savinnsk/prototype_bot_whatsapp/internal/entity"

func CreateUser(user *entity.User) {

	connection := Init()
	err := connection.db.Create(&user).Error
	if err != nil {
		println("Error to create user")
	}
}

func FindUserByPhone(user *entity.User) {
	connection := Init()
	connection.db.Where("phone = ?", user.Phone).First(&user)
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
