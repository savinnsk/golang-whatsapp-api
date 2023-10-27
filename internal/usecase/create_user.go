package usecase

import "github.com/savinnsk/prototype_bot_whatsapp/internal/entity"


type UserDTO struct {
	Name string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone" binding:"required"`
}


type CreateUserUseCase struct {
	UserRepository  any //UserRepository events
	UserCreated any //UserCreated events
	EventDispatcher any //EventDispatcher events
}


func (c *CreateUserUseCase) Execute(userDTO UserDTO) {
	user := entity.User{
		Name: userDTO.Name,
		PhoneNumber: userDTO.PhoneNumber,
		IsActive: true,
		Role: "user",
	}


	// if err := c.UserRepository.Save(user); err != nil {
	// 	return UserDTO{}, err
	// }
	// c.UserCreated.SetPayload(user)
	// c.UserCreated.Dispatch(c.OrderCreated)
	// return user, nil
	println(user)
	
	// user := c.UserRepository.Create(userDTO)

	// c.EventDispatcher.Dispatch(user)
}