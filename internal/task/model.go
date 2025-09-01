package task

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	UserID uint32
	Title  string
	IsDone bool
}
