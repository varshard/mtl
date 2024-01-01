package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/varshard/mtl/api/handlers"
	"github.com/varshard/mtl/api/middlewares"
	"github.com/varshard/mtl/api/routes"
	"github.com/varshard/mtl/infrastructure/config"
	"github.com/varshard/mtl/infrastructure/database"
	"github.com/varshard/mtl/infrastructure/database/repository"
	"gorm.io/gorm"
	"net/http"
)

type Server struct {
}

func (s Server) Start(conf *config.Config) {
	db, err := database.InitDB(&conf.DBConfig)
	if err != nil {
		panic(fmt.Sprintf("fail to connect to the database: %s", err.Error()))
	}
	r := s.InitRoutes(db, conf)
	if err := http.ListenAndServe(conf.Port, r); err != nil {
		panic(fmt.Sprintf("fail to start server: %s", err.Error()))
	}
}

func (s Server) InitRoutes(db *gorm.DB, conf *config.Config) *chi.Mux {
	userRepo := repository.UserRepository{DB: db}
	voteItemRepository := repository.ItemRepository{DB: db}

	authHandler := &handlers.AuthHandler{UserRepository: userRepo, Config: conf}
	authMiddleware := middlewares.NewAuthenticationMiddleware(conf.Secret, userRepo)

	r := chi.NewRouter()
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello"))
	})
	r.Post("/login", authHandler.Login)
	r.Mount("/vote_items", routes.MakeVoteItemsRoutes(voteItemRepository, authMiddleware))
	return r
}
