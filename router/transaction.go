package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rammyblog/go-product-subscriptions/controller"
	"github.com/rammyblog/go-product-subscriptions/middleware"
)

func TransactionRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{"message": "welcome to transactions"})
	})
	r.Post("/webhook", controller.TransactionWebhook)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware)
		r.Post("/init-customer", controller.InitializeCustomerTransaction)
	})

	return r
}
