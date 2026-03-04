package user

import "time"

type UserStruct struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	Email string 
	Password string 
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (UserStruct) TableName() string {
    return "user_structs" 
}