package main

import (
	"github.com/varshard/mtl/infrastructure/config"
	"github.com/varshard/mtl/infrastructure/server"
)

func main() {
	conf := config.ReadEnv()
	s := server.Server{}
	s.Start(&conf)
}
