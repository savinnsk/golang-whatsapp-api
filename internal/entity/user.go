package entity

type User struct {
	Id    int    `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Phone string `gorm:"not null"`
	Role  string `gorm:"not null"`
}

func NewUser(name string, phone string, role string) *User {
	return &User{
		Name:  name,
		Phone: phone,
		Role:  role,
	}
}
