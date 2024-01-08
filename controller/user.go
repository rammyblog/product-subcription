package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rammyblog/go-product-subscriptions/database"
	"github.com/rammyblog/go-product-subscriptions/helper"
	"github.com/rammyblog/go-product-subscriptions/middleware"
	"github.com/rammyblog/go-product-subscriptions/models"
	"github.com/rammyblog/go-product-subscriptions/response"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}
	input := &helper.InputRequest{
		Data: user,
	}
	if err := input.Bind(r); err != nil {
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}

	var existingUser models.User

	database.DB.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != 0 {
		render.Render(w, r, response.ErrInvalidRequest(errors.New("email already exist")))
		return
	}

	encpw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}

	user.Password = string(encpw)
	// create user
	result := database.DB.Create(&user)

	if result.Error != nil {
		render.Render(w, r, response.ErrInvalidRequest(result.Error))
		return
	}

	render.Render(w, r, response.Response(http.StatusCreated, user))
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var input LoginRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Fatal(err)
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}
	var user models.User

	bindData := &helper.InputRequest{
		Data: input,
	}

	if err := bindData.Bind(r); err != nil {
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}

	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		render.Render(w, r, response.ErrInvalidRequest(errors.New("record not found")))
		return
	}
	// compare password
	validPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) == nil

	if !validPassword {
		render.Render(w, r, response.ErrInvalidRequest(errors.New("email or password incorrect")))
		return
	}

	token, err := middleware.CreateJwtToken(user)
	if err != nil {
		log.Fatal(err)
		render.Render(w, r, response.ErrInvalidRequest(errors.New("email or password incorrect")))
		return
	}
	render.Render(w, r, response.Response(http.StatusOK, map[string]string{"token": token}))
}
