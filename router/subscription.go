package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rammyblog/go-product-subscriptions/controller"
	"github.com/rammyblog/go-product-subscriptions/middleware"
)

func SubscriptionRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"message": "welcome to subscriptions"})
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware)
		r.Post("/", controller.CreateSubscription)
		r.Get("/", controller.GetSubscriptions)
		r.Get("/{id}", controller.GetSubscription)
		// r.Delete("/{id}", controller.DeleteSubscription)
	})

	return r

}
