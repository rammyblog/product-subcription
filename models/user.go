package models

type User struct {
	Model
	Name     string `json:"name"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;" json:"email"`
	Password string `gorm:"not null" json:"-"`
}
