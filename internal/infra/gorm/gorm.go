package infra

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Schedule struct {
	Id        int    `gorm:"primaryKey"`
	Time      string `gorm:"not null"`
	Date      *string
	Available bool `gorm:"not null"`
	Disabled  bool `gorm:"not null"`
	gorm.Model
}

type User struct {
	Id         int        `gorm:"primaryKey"`
	Name       string     `gorm:"not null"`
	Role       float64    `gorm:"not null"`
	Schedules  []Schedule `gorm:"many2many:user_schedules"`
	gorm.Model            // to create created_at and updated_at and deleted_at
}

type Connection struct {
	db *gorm.DB
}

func Init() *Connection {

	db, err := gorm.Open(sqlite.Open("file:data.db?_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	connection := Connection{db: db}
	connection.db.AutoMigrate(&Schedule{}, &User{})

	println("Done")

	return &connection
}
