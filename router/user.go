package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rammyblog/go-product-subscriptions/controller"
	"github.com/rammyblog/go-product-subscriptions/middleware"
)

func UserRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"message": "welcome to users"})
	})

	r.Post("/", controller.CreateUser)
	r.Post("/login", controller.LoginUser)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware)
		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]string{"message": "me"})
		})
	})

	return r
}
