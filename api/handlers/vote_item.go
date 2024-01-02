package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/varshard/mtl/api/handlers/responses"
	"github.com/varshard/mtl/api/middlewares"
	"github.com/varshard/mtl/domain/vote"
	"github.com/varshard/mtl/infrastructure/database"
	"github.com/varshard/mtl/infrastructure/database/repository"
	xErr "github.com/varshard/mtl/infrastructure/errors"
	"github.com/varshard/mtl/infrastructure/rest"
	"net/http"
	"strconv"
)

type VoteItemHandler struct {
	UserRepository     repository.UserRepository
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

type VoteItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (v VoteItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	payload := VoteItemRequest{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		rest.ServeJSON(http.StatusBadRequest, w, &responses.ErrorResponse{Error: err.Error()})
		return
	}

	ctx := r.Context()
	claims, _ := ctx.Value(middlewares.AuthContext).(jwt.StandardClaims)

	usr, err := v.UserRepository.FindUser(claims.Subject)

	if err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: responses.ErrInternalServerError})
		return
	}

	item, err := v.VoteItemRepository.Create(database.VoteItem{
		Name:        payload.Name,
		Description: payload.Description,
		CreatedBy:   usr.ID,
	})

	if err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: responses.ErrInternalServerError})
		return
	}

	rest.ServeJSON(http.StatusOK, w, responses.NewDataResponse(item))
}

func (v VoteItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	payload := VoteItemRequest{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		rest.ServeJSON(http.StatusBadRequest, w, &responses.ErrorResponse{Error: "invalid payload"})
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		rest.ServeJSON(http.StatusBadRequest, w, &responses.ErrorResponse{Error: "id isn't a number"})
		return
	}

	err = v.VoteItemRepository.Update(uint(id), vote.UpdateVoteItem{
		Name:        payload.Name,
		Description: payload.Description,
	})

	if errors.As(err, &xErr.ErrNotFound{}) {
		rest.ServeJSON(http.StatusNotFound, w, &responses.ErrorResponse{Error: "vote item not found"})
		return
	} else if err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: responses.ErrInternalServerError})
		return
	}

	rest.ServeJSON(http.StatusOK, w, responses.NewSuccessResponse(true))
}

func (v VoteItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		rest.ServeJSON(http.StatusBadRequest, w, &responses.ErrorResponse{Error: "id isn't a number"})
		return
	}

	itemID := uint(id)

	ok, err := v.VoteItemRepository.Removable(itemID)
	if errors.As(err, &xErr.ErrNotFound{}) {
		rest.ServeJSON(http.StatusNotFound, w, &responses.ErrorResponse{Error: "vote item not found"})
		return
	} else if err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: responses.ErrInternalServerError})
		return
	}

	if !ok {
		rest.ServeJSON(http.StatusBadRequest, w, &responses.ErrorResponse{Error: "the item has been voted"})
		return
	}

	if err := v.VoteItemRepository.Remove(itemID); err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: "fail to remove an item"})
		return
	}

	rest.ServeJSON(http.StatusOK, w, responses.NewSuccessResponse(true))
}

func (v VoteItemHandler) ResetItems(w http.ResponseWriter, r *http.Request) {
	if err := v.VoteItemRepository.ResetItems(); err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: "fail to remove an item"})
		return
	}
	rest.ServeJSON(http.StatusOK, w, responses.NewSuccessResponse(true))
}
