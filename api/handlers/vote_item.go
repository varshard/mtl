package handlers

import (
	"github.com/varshard/mtl/api/handlers/responses"
	"github.com/varshard/mtl/infrastructure/database/repository"
	"github.com/varshard/mtl/infrastructure/rest"
	"net/http"
)

type VoteItemHandler struct {
	VoteItemRepository repository.ItemRepository
}

func (v VoteItemHandler) GetVoteItems(w http.ResponseWriter, r *http.Request) {
	items, err := v.VoteItemRepository.GetItems()
	if err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: responses.ErrInternalServerError})
		return
	}

	rest.ServeJSON(http.StatusOK, w, responses.NewDataResponse(items))
}
