package main

import (
	"github.com/varshard/mtl/infrastructure/config"
)

func main() {
	conf := config.ReadEnv()
	s := Server{}
	s.Start(&conf)
}
