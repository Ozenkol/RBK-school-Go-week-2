package router

import (
	"net/http"
)

func NewRouter(h *handler.Handler) http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /health", h.Health)
    mux.HandleFunc("GET /api/v1/users/{id}", h.GetUser)
    mux.HandleFunc("POST /api/v1/users", h.CreateUser)

    return mux
}