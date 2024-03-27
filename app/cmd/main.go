package main

import (
	"github.com/Mooonsheen/lamoda_tech/app/internal/server"
	"github.com/Mooonsheen/lamoda_tech/app/internal/server/config"
)

func main() {
	cfg := new(config.Config)
	cfg.Read()
	server := server.NewServer(cfg)
	server.Run()
}
