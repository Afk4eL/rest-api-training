package models

import (
	"time"
)

type UserEntity struct {
	Id        uint   `json:"Id" gorm:"primarykey"`
	Username  string `json:"username" gorm:"notnull" validate:"required"`
	Email     string `json:"email" gorm:"unique;notnull" validate:"required"`
	Password  string `json:"password" gorm:"notnull" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserEntity) TableName() string {
	return "users"
}
