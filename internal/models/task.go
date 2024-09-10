package models

import "time"

type TaskEntity struct {
	Id        uint   `json:"id" gorm:"primarykey"`
	UserId    uint   `json:"user_id" gorm:"notnull;column:user_id"`
	TaskId    uint   `json:"task_id" gorm:"notnull;column:task_id"`
	Title     string `json:"title" gorm:"notnull" validate:"required"`
	Data      string `json:"data"  validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (TaskEntity) TableName() string {
	return "tasks"
}
