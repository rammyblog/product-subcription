package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/rammyblog/go-paystack"
	"github.com/rammyblog/go-product-subscriptions/config"
	"github.com/rammyblog/go-product-subscriptions/middleware"
	"github.com/rammyblog/go-product-subscriptions/models"
	"github.com/rammyblog/go-product-subscriptions/response"
)

// CreateSubscription handles the creation of a new subscription.
// It expects a JSON payload containing the subscription details in the request body.
// The user must be authenticated with a valid token in the Authorization header.
// It retrieves the current user and the product associated with the subscription from the database.
// Then, it creates a subscription using the Paystack API and saves it to the database.
// Finally, it returns the created subscription as a JSON response.
func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var subscription models.Subscription
	var currentUser models.User
	var product models.Product

	err := json.NewDecoder(r.Body).Decode(&subscription)
	if err != nil {
		render.JSON(w, r, response.ErrInvalidRequest(err))
		return
	}

	userID, err := middleware.GetUserIdFromToken(r.Header.Get("Authorization"))
	if err != nil {
		render.Render(w, r, response.ErrInvalidRequest(errors.New("unauthorized")))
		return
	}
	if err := config.GlobalConfig.DB.Where("id = ?", userID).First(&currentUser).Error; err != nil {
		render.Render(w, r, response.ErrInvalidRequest(errors.New("record not found")))
		return
	}
	if err := config.GlobalConfig.DB.Where("id = ?", subscription.ProductID).First(&product).Error; err != nil {
		render.Render(w, r, response.ErrInvalidRequest(errors.New("record not found")))
		return
	}

	subscription.UserID = currentUser.ID
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payStackSubscription, err := config.GlobalConfig.PaystackClient.Subscription.CreateSubscription(ctx, &paystack.CreateSubscriptionRequest{
		Customer: currentUser.CustomerCode,
		Plan:     product.PlanCode,
	})

	if err != nil {
		jsonErr, _ := json.Marshal(err)
		var errMap map[string]interface{}

		_ = json.Unmarshal(jsonErr, &errMap)
		errMessage := errMap["details"].(map[string]interface{})["message"].(string)
		if errMessage == "This subscription is already in place." {
			// TODO: update the subscription status in the database
			render.JSON(w, r, response.ErrInvalidRequest(errors.New("this subscription is already in place")))
			return
		}
		render.JSON(w, r, response.ErrInvalidRequest(err))
		return
	}

	subscription.SubscriptionCode = payStackSubscription.SubscriptionCode
	subscription.Status = payStackSubscription.Status

	if err := config.GlobalConfig.DB.Create(&subscription).Error; err != nil {
		render.JSON(w, r, response.ErrInvalidRequest(err))
		return
	}

	render.JSON(w, r, response.Response(http.StatusOK, subscription))
}
