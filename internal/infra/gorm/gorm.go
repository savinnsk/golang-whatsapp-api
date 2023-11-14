package infra

import (
	entity "github.com/savinnsk/prototype_bot_whatsapp/internal/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Connection struct {
	db *gorm.DB
}

func Init() *Connection {

	db, err := gorm.Open(sqlite.Open("file:data.db?_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	connection := Connection{db: db}
	connection.db.AutoMigrate(&entity.Schedule{}, &entity.User{}, &entity.UserSchedule{})

	println("Connection with database ok")

	return &connection
}
