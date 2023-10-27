package entity


type User struct {
	Id int
	Name string
	PhoneNumber string
	Role string
	IsActive bool
}


func NewUser(name string, phoneNumber string,role string , isActive bool) *User {
	return &User{
		Name: name,
		PhoneNumber: phoneNumber,
		Role: role,
		IsActive: isActive,
	}
}