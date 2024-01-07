package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rammyblog/go-product-subscriptions/controller"
)

func UserRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"message": "welcome to users"})
	})

	r.Post("/", controller.CreateUser)

	return r
}
