package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/varshard/mtl/api/handlers/responses"
	"github.com/varshard/mtl/api/middlewares"
	"github.com/varshard/mtl/infrastructure/database/repository"
	"github.com/varshard/mtl/infrastructure/rest"
	"net/http"
	"strconv"
)

const ErrInvalidItemID = "id isn't a number"

type VoteHandler struct {
	VoteRepository     repository.VoteRepository
	VoteItemRepository repository.ItemRepository
	UserRepository     repository.UserRepository
}

func (v VoteHandler) Vote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		rest.ServeJSON(http.StatusBadRequest, w, &responses.ErrorResponse{Error: ErrInvalidItemID})
		return
	}

	ctx := r.Context()
	claims, _ := ctx.Value(middlewares.AuthContext).(jwt.StandardClaims)

	usr, err := v.UserRepository.FindUser(claims.Subject)

	if err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: responses.ErrInternalServerError})
		return
	}

	itemID := uint(id)
	exist, err := v.VoteItemRepository.Exist(itemID)
	if err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: responses.ErrInternalServerError})
		return
	}
	if !exist {
		rest.ServeJSON(http.StatusNotFound, w, &responses.ErrorResponse{Error: "vote item not found"})
		return
	}
	ok, err := v.VoteRepository.IsVoteable(usr.ID)
	if err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: responses.ErrInternalServerError})
		return
	}
	if !ok {
		rest.ServeJSON(http.StatusBadRequest, w, &responses.ErrorResponse{Error: "you already casted your vote"})
		return
	}

	if err := v.VoteRepository.Vote(itemID, usr.ID); err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: responses.ErrInternalServerError})
		return
	}
	rest.ServeJSON(http.StatusOK, w, responses.NewSuccessResponse(true))
}

func (v VoteHandler) ClearVotes(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		rest.ServeJSON(http.StatusBadRequest, w, &responses.ErrorResponse{Error: ErrInvalidItemID})
		return
	}

	itemID := uint(id)

	if err := v.VoteRepository.ClearVote(itemID); err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: responses.ErrInternalServerError})
		return
	}
	rest.ServeJSON(http.StatusOK, w, responses.NewSuccessResponse(true))
}
