package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rammyblog/go-product-subscriptions/helper"
	"github.com/rammyblog/go-product-subscriptions/models"
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
	render.JSON(w, r, input)

}
