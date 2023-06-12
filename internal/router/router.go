package router

import (
	"github.com/MWT-proger/compressfile/internal/handlers"
	"github.com/go-chi/chi"
)

// Router() Перенаправляет запросы на необходимые хендлеры
func Router(h *handlers.APIHandler) *chi.Mux {

	r := chi.NewRouter()
	r.Get("/", h.TransformImage)

	return r
}
