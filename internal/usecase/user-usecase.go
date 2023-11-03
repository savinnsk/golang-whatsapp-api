package usecase

import (
	"github.com/satori/go.uuid"
	dto "github.com/savinnsk/prototype_bot_whatsapp/internal/domain/dto"
	entity "github.com/savinnsk/prototype_bot_whatsapp/internal/entity"
	gorm "github.com/savinnsk/prototype_bot_whatsapp/internal/infra/gorm"
)

func CreateUserAndReturn(data dto.CreateUserDto) (*entity.User, error) {
	user := entity.NewUser(data.Name, data.Phone, data.Role, uuid.NewV4().String())
	err := gorm.CreateUser(user)
	if err != nil {
		println("Error to create user")
		return nil, err
	}

	userFounded, err := gorm.FindUserByPhone(data.Phone)
	if err != nil {
		println("Error to find user")
		return nil, err
	}

	return userFounded, nil
}
