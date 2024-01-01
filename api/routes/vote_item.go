package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/varshard/mtl/api/handlers"
	"github.com/varshard/mtl/infrastructure/database/repository"
	"net/http"
)

func MakeVoteItemsRoutes(voteItemRepository repository.ItemRepository, middlewares ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middlewares...)

	handler := handlers.VoteItemHandler{VoteItemRepository: voteItemRepository}

	r.Get("/", handler.GetVoteItems)

	return r
}
