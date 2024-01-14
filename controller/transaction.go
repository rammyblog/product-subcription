package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/rammyblog/go-paystack"
	"github.com/rammyblog/go-product-subscriptions/config"
	"github.com/rammyblog/go-product-subscriptions/helper"
	"github.com/rammyblog/go-product-subscriptions/middleware"
	"github.com/rammyblog/go-product-subscriptions/models"
	"github.com/rammyblog/go-product-subscriptions/response"
)

func InitializeCustomerTransaction(w http.ResponseWriter, r *http.Request) {

	// get the user id from the token
	userID, err := middleware.GetUserIdFromToken(r.Header.Get("Authorization"))
	if err != nil {
		render.Render(w, r, response.ErrInvalidRequest(errors.New("unauthorized")))
		return
	}

	var currentUser models.User

	// get the user from the database
	if err := config.GlobalConfig.DB.Where("id = ?", userID).First(&currentUser).Error; err != nil {
		render.Render(w, r, response.ErrInvalidRequest(errors.New("record not found")))
		return
	}

	// create a transaction in paystack for customer
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	randomAlphabet := helper.GenerateRandomAlphabet()

	transaction, err := config.GlobalConfig.PaystackClient.Transaction.Initialize(ctx, &paystack.TransactionRequest{
		Reference: randomAlphabet,
		Amount:    1000,
		Email:     currentUser.Email,
	})

	if err != nil {
		log.Fatal(err)
		render.Render(w, r, response.ErrInvalidRequest(errors.New("unable to create transaction")))
		return
	}

	render.Render(w, r, response.Response(http.StatusOK, transaction))

}
