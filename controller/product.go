package controller

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rammyblog/go-product-subscriptions/database"
	"github.com/rammyblog/go-product-subscriptions/models"
	"github.com/rammyblog/go-product-subscriptions/response"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	database.DB.Find(&products)
	render.JSON(w, r, response.Response(http.StatusOK, products))
}
