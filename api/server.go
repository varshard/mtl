package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/varshard/mtl/api/handlers"
	"github.com/varshard/mtl/api/middlewares"
	"github.com/varshard/mtl/infrastructure/config"
	"github.com/varshard/mtl/infrastructure/database"
	"github.com/varshard/mtl/infrastructure/database/repository"
	"gorm.io/gorm"
	"net/http"
)

type MTLServer struct {
}

func (s MTLServer) Start(conf *config.Config) {
	db, err := database.InitDB(&conf.DBConfig)
	if err != nil {
		panic(fmt.Sprintf("fail to connect to the database: %s", err.Error()))
	}
	r := s.InitRoutes(db, conf)
	if err := http.ListenAndServe(conf.Port, r); err != nil {
		panic(fmt.Sprintf("fail to start server: %s", err.Error()))
	}
}

func (s MTLServer) InitRoutes(db *gorm.DB, conf *config.Config) *chi.Mux {
	userRepo := repository.UserRepository{DB: db}
	voteItemRepository := repository.ItemRepository{DB: db}
	voteRepository := repository.VoteRepository{DB: db}

	authHandler := &handlers.AuthHandler{UserRepository: userRepo, Config: conf}
	authMiddleware := middlewares.NewAuthenticationMiddleware(conf.Secret, userRepo)

	voteItemHandler := handlers.VoteItemHandler{UserRepository: userRepo, VoteItemRepository: voteItemRepository}
	voteHandler := handlers.VoteHandler{VoteRepository: voteRepository, UserRepository: userRepo, VoteItemRepository: voteItemRepository}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello"))
	})
	r.Post("/login", authHandler.Login)

	voteRouter := chi.NewRouter()
	voteRouter.Use(authMiddleware)
	voteRouter.Get("/", voteItemHandler.GetVoteItems)
	voteRouter.Post("/", voteItemHandler.CreateItem)
	voteRouter.Post("/reset", voteItemHandler.ResetItems)
	voteRouter.Put("/{id}", voteItemHandler.Update)
	voteRouter.Delete("/{id}", voteItemHandler.Delete)
	voteRouter.Post("/{id}/vote", voteHandler.Vote)
	voteRouter.Post("/{id}/reset", voteHandler.ClearVotes)

	r.Mount("/vote_items", voteRouter)
	return r
}
