package server

import (
	"context"
	"fmt"

	config "github.com/Mooonsheen/lamoda_tech/app/internal/server/config"
	"github.com/Mooonsheen/lamoda_tech/app/internal/storage"
	configdb "github.com/Mooonsheen/lamoda_tech/app/internal/storage/config"
	"github.com/Mooonsheen/lamoda_tech/app/internal/storage/interfaces"
	"github.com/Mooonsheen/lamoda_tech/app/internal/storage/postgresql"

	// "github.com/Mooonsheen/lamoda_tech/app/internal/storage/postgresql"
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

// var InternalCache cache.Cache // разобраться с интерфейсом

func (s *Server) Run() {
	cfgDb := new(configdb.ConfigDb)
	cfgDb.Read()
	var err error
	Pool, err := storage.NewStorageClient(context.TODO(), cfgDb)
	if err != nil {
		fmt.Println(err)
	}
	defer Pool.Close()
	PgClient = postgresql.NewDatadase(Pool, Pool)

	r := gin.Default()
	r.POST("/reservation", s.handleReservation) // Создает бронь
	r.PATCH("/reservation", s.handleReservation)
	r.DELETE("/reservation", s.handleReservation) // Удаляет бронь (товар забран со склада или истек ttl)
	// r.GET("/store/{id}", s.handleGetStoreRemains)

	serverPath := s.cfg.Host + s.cfg.Port
	r.Run(serverPath)
}
