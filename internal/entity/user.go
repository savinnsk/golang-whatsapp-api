package entity

type User struct {
	Id    string `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Phone string `gorm:"not null"`
	Role  string `gorm:"not null"`
}

func NewUser(name string, phone string, role string, id string) *User {
	return &User{
		Id:    id,
		Name:  name,
		Phone: phone,
		Role:  role,
	}
}
