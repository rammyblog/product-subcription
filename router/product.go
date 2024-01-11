package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/rammyblog/go-product-subscriptions/controller"
)

func ProductRoutes() chi.Router {

	r := chi.NewRouter()

	r.Get("/", controller.GetAllProducts)

	return r

}
