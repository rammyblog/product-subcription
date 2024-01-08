package models

type User struct {
	Model
	Name     string `json:"name" validate:"required"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;" json:"email" validate:"required,email"`
	Password string `gorm:"not null" json:"password" validate:"required"`
}
