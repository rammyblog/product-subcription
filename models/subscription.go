package models

type Subscription struct {
	Model
	ProductID        uint   `json:"product_id" validate:"required"`
	UserID           uint   `json:"user_id" validate:"required"`
	SubscriptionCode string `json:"subscription_code" validate:"required"`
	Status           string `json:"status" validate:"required"`
}
