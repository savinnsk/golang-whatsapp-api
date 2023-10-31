package entity

type User struct {
	Id        int         `gorm:"primaryKey"`
	Name      string      `gorm:"not null"`
	Phone     string      `gorm:"not null"`
	Role      string      `gorm:"not null"`
	Schedules *[]Schedule `gorm:"many2many:user_schedules"` // to create created_at and updated_at and deleted_at
}

func NewUser(name string, phone string, role string) *User {
	return &User{
		Name:  name,
		Phone: phone,
		Role:  role,
	}
}
