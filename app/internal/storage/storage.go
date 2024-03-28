package storage

import (
	"context"
	"fmt"
	"time"

	configdb "github.com/Mooonsheen/lamoda_tech/app/internal/storage/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	maxAttempts = 3
)

func NewStorageClient(ctx context.Context, cfg *configdb.ConfigDb) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 3*time.Second)

	if err != nil {
		fmt.Printf("error DoWithTries postgresql, dsn: %s", dsn)
	}

	return pool, nil
}

func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}

		return nil
	}

	return
}
