package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/varshard/mtl/api/handlers"
	"github.com/varshard/mtl/api/middlewares"
	"github.com/varshard/mtl/infrastructure/config"
	"github.com/varshard/mtl/infrastructure/database/repository"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type MTLServer struct {
	Server           *http.Server
	DB               *gorm.DB
	shutdownFinished chan struct{}
}

func (s *MTLServer) Start(conf *config.Config) {
	if s.shutdownFinished == nil {
		s.shutdownFinished = make(chan struct{})
	}
	r := s.InitRoutes(s.DB, conf)

	s.Server = &http.Server{Addr: conf.Port, Handler: r}

	err := s.Server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server has been shutdown")
	} else if err != nil {
		fmt.Printf("fail to start the server %s\n", err.Error())
	}
	<-s.shutdownFinished
	fmt.Println("server finished shutting down")
}

func (s *MTLServer) InitRoutes(db *gorm.DB, conf *config.Config) *chi.Mux {
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

func (s *MTLServer) GracefulExit() {
	waiter := make(chan os.Signal, 1)
	signal.Notify(waiter, syscall.SIGTERM, os.Interrupt)
	<-waiter

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("closing db connection")
	dbInstance, _ := s.DB.DB()
	dbInstance.Close()
	fmt.Println("close db connection")

	if err := s.Server.Shutdown(ctx); err != nil {
		fmt.Printf("fail to shut down the server with error %s\n", err.Error())
	}
	close(s.shutdownFinished)
}
