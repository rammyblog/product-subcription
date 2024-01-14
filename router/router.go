package router

import (
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type CustomRenderer struct {
	templates *template.Template
}

func (cr *CustomRenderer) Render(w http.ResponseWriter, r *http.Request, name string, data interface{}, _ chi.Route) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return cr.templates.ExecuteTemplate(w, name, data)
}

func (cr *CustomRenderer) Bind(bind func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := bind(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
}

func Init() *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", UserRoutes())
		r.Mount("/products", ProductRoutes())
		r.Mount("/transactions", TransactionRouter())
		r.Mount("/subscriptions", SubscriptionRouter())

	})

	return r
}
