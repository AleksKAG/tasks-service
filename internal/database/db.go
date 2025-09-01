package database

import (
	"github.com/AleksKAG/tasks-service/internal/task"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("tasks.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	DB.AutoMigrate(&task.Task{})
}
