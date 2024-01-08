package controller

import (
	"net/http"

	"github.com/go-chi/render"
)

type SuccessResponse struct {
	HTTPStatusCode int         `json:"-"` // http response status code
	Data           interface{} `json:"data"`
	Status         bool        `json:"status"`
}

func (e *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func Response(httpCode int, data interface{}) render.Renderer {
	return &SuccessResponse{
		Status:         true,
		HTTPStatusCode: httpCode,
		Data:           data,
	}
}
