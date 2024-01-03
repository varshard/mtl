package main

import (
	"fmt"
	"github.com/varshard/mtl/api"
	"github.com/varshard/mtl/infrastructure/config"
	"github.com/varshard/mtl/infrastructure/database"
)

func main() {
	conf := config.ReadEnv()
	db, err := database.InitDB(&conf.DBConfig)
	if err != nil {
		panic(fmt.Sprintf("fail to connect to the database: %s", err.Error()))
	}
	s := api.MTLServer{DB: db}
	go s.GracefulExit()
	s.Start(&conf)
}
