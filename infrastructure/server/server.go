package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/varshard/mtl/handlers"
	"github.com/varshard/mtl/infrastructure/config"
	"github.com/varshard/mtl/infrastructure/database"
	"gorm.io/gorm"
	"net/http"
)

type Server struct {
}

func (s Server) Start(conf *config.Config) {
	db, err := database.InitDB(conf)
	if err != nil {
		panic(fmt.Sprintf("fail to connect to the database: %s", err.Error()))
	}
	r := s.InitRoutes(db, conf)
	if err := http.ListenAndServe(conf.Port, r); err != nil {
		panic(fmt.Sprintf("fail to start server: %s", err.Error()))
	}
}

func (s Server) InitRoutes(db *gorm.DB, conf *config.Config) *chi.Mux {
	authCtrl := &handlers.AuthHandler{DB: db, Config: conf}

	r := chi.NewRouter()
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello"))
	})
	r.Post("/login", authCtrl.Login)
	return r
}
