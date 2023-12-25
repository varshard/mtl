package handlers

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/varshard/mtl/domain/user"
	"github.com/varshard/mtl/infrastructure/config"
	"github.com/varshard/mtl/infrastructure/rest"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

const (
	ErrInvalidCredentials  = "invalid username or password"
	ErrInternalServerError = "internal server error"
)

type AuthHandler struct {
	DB     *gorm.DB
	Config *config.Config
}

type (
	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	LoginResponse struct {
		Token string `json:"token"`
	}
	ErrorResponse struct {
		Error string `json:"error"`
	}
)

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	payload := Login{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		rest.ServeJSON(http.StatusBadRequest, w, &ErrorResponse{Error: err.Error()})
		return
	}

	u := user.User{}
	err := a.DB.Table("user").Select("id, name, password").
		Where("name = ?", payload.Username).Limit(1).First(&u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		rest.ServeJSON(http.StatusUnauthorized, w, &ErrorResponse{Error: ErrInvalidCredentials})
		return
	} else if err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &ErrorResponse{Error: ErrInternalServerError})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(payload.Password)); err != nil {
		rest.ServeJSON(http.StatusUnauthorized, w, &ErrorResponse{Error: ErrInvalidCredentials})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject: u.Name,
	})
	tokenString, err := token.SignedString([]byte(a.Config.Secret))
	if err != nil {
		rest.ServeJSON(http.StatusInternalServerError, w, &ErrorResponse{Error: ErrInternalServerError})
		return
	}
	rest.ServeJSON(http.StatusOK, w, &LoginResponse{Token: tokenString})
	return
}
