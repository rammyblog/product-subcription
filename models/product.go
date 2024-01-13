package models

type Product struct {
	Model
	Name        string  `json:"name" validate:"required"`
	Description string  `gorm:"type:varchar(200);" json:"description" validate:"required"`
	Price       float64 `gorm:"not null" json:"price" validate:"required"`
	Duration    string  `gorm:"not null" json:"duration" validate:"required,oneof=month yearly"`
	PlanCode    string  `gorm:"not null" json:"plan_code" validate:"required"`
}
