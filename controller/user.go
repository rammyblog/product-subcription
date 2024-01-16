package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rammyblog/go-product-subscriptions/config"
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

type UserResponse struct {
	ID           uint      `json:"id"`
	Email        string    `json:"email"`
	CustomerCode string    `json:"customer_code"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
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

	config.GlobalConfig.DB.Where("email = ?", user.Email).First(&existingUser)
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
	result := config.GlobalConfig.DB.Create(&user)

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

	if err := config.GlobalConfig.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var user models.User
	if err := config.GlobalConfig.DB.Where("id = ?", id).First(&user).Error; err != nil {
		render.Render(w, r, response.ErrInvalidRequest(errors.New("record not found")))
		return
	}
	render.Render(w, r, response.Response(http.StatusOK, user))
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIdFromToken(r.Header.Get("Authorization"))
	if err != nil {
		render.Render(w, r, response.ErrInvalidRequest(errors.New("unauthorized")))
		return
	}
	var currentUser models.User
	if err := config.GlobalConfig.DB.Where("id = ?", userID).First(&currentUser).Error; err != nil {
		render.Render(w, r, response.ErrInvalidRequest(errors.New("record not found")))
		return
	}

	responseUser := UserResponse{
		ID:           currentUser.ID,
		Email:        currentUser.Email,
		UpdatedAt:    currentUser.UpdatedAt,
		CreatedAt:    currentUser.CreatedAt,
		CustomerCode: currentUser.CustomerCode,
	}

	render.Render(w, r, response.Response(http.StatusOK, responseUser))
}
