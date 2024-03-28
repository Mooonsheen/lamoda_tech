package server

import (
	"context"
	"fmt"

	config "github.com/Mooonsheen/lamoda_tech/app/internal/server/config"
	"github.com/Mooonsheen/lamoda_tech/app/internal/storage"
	configdb "github.com/Mooonsheen/lamoda_tech/app/internal/storage/config"
	"github.com/Mooonsheen/lamoda_tech/app/internal/storage/interfaces"
	"github.com/Mooonsheen/lamoda_tech/app/internal/storage/postgresql"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

var Pool *pgxpool.Pool
var PgClient interfaces.Storage

func (s *Server) Run() {
	cfgDb := new(configdb.ConfigDb)
	cfgDb.Read()
	var err error
	Pool, err := storage.NewStorageClient(context.TODO(), cfgDb)
	if err != nil {
		fmt.Println(err)
	}
	defer Pool.Close()
	PgClient = postgresql.NewDatadase(Pool)

	r := gin.Default()
	r.POST("/reservation", s.handleReservation)
	r.PATCH("/reservation", s.handleReservation)
	r.DELETE("/reservation", s.handleReservation)
	r.GET("/store", s.handleGetStoreRemains)

	serverPath := s.cfg.Host + s.cfg.Port
	r.Run(serverPath)
}
