package models

type Subscription struct {
	Model
	ProductID        uint   `json:"product_id" validate:"required"`
	UserID           uint   `json:"user_id"`
	SubscriptionCode string `json:"subscription_code"`
	Status           string `json:"status"`
}
