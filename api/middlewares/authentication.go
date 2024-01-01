package middlewares

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/varshard/mtl/api/handlers/responses"
	"github.com/varshard/mtl/infrastructure/database/repository"
	"github.com/varshard/mtl/infrastructure/rest"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

const AuthContext = "auth"

func NewAuthenticationMiddleware(secret string, repository repository.UserRepository) func(http.Handler) http.Handler {
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			claims, valid := authenticateClaims(secret, r.Header.Get("Authorization"))
			if !valid {
				rest.ServeJSON(http.StatusUnauthorized, w, &responses.ErrorResponse{Error: responses.ErrUnauthorized})
				return
			}
			username := claims.Subject

			_, err := repository.FindUser(username)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				rest.ServeJSON(http.StatusUnauthorized, w, &responses.ErrorResponse{Error: responses.ErrUnauthorized})
				return
			} else if err != nil {
				rest.ServeJSON(http.StatusInternalServerError, w, &responses.ErrorResponse{Error: responses.ErrInternalServerError})
				return
			}

			ctx = context.WithValue(ctx, AuthContext, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	return middleware
}

func authenticateClaims(secret, authHeader string) (jwt.StandardClaims, bool) {
	authSplitted := strings.Split(authHeader, " ")
	if len(authSplitted) < 2 {
		return jwt.StandardClaims{}, false
	}
	if !strings.EqualFold(authSplitted[0], "Bearer") {
		return jwt.StandardClaims{}, false
	}

	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(authSplitted[1], &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return jwt.StandardClaims{}, false
	}
	if token == nil {
		return jwt.StandardClaims{}, false
	}

	return claims, true
}
