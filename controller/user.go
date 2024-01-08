package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rammyblog/go-product-subscriptions/database"
	"github.com/rammyblog/go-product-subscriptions/helper"
	"github.com/rammyblog/go-product-subscriptions/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	input := &helper.InputRequest{
		Data: user,
	}
	if err := input.Bind(r); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var existingUser models.User

	database.DB.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != 0 {
		render.Render(w, r, ErrInvalidRequest(errors.New("email already exist")))
		return
	}

	encpw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user.Password = string(encpw)
	// create user
	result := database.DB.Create(&user)

	if result.Error != nil {
		render.Render(w, r, ErrInvalidRequest(result.Error))
		return
	}

	render.Render(w, r, Response(http.StatusCreated, user))

}
