package main

import (
	"github.com/varshard/mtl/api"
	"github.com/varshard/mtl/infrastructure/config"
)

func main() {
	conf := config.ReadEnv()
	s := api.MTLServer{}
	s.Start(&conf)
}
