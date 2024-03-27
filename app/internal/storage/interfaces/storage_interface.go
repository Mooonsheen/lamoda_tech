package interfaces

import (
	"context"

	"github.com/Mooonsheen/lamoda_tech/app/internal/models"
)

type Storage interface {
	CreateReservation(ctx context.Context, msg models.ReserveRequestMessage, uuid string) (reservedItems models.ResponseMessageAfterCreate, err error)
	СonfirmReservation(ctx context.Context, uuid string) (status string, err error)
	DeleteReservation(ctx context.Context, uuid string) (status string, err error)
	GetStoreRemains(ctx context.Context, store_id string) (remains map[string]int, err error)
	GetStoresAvailability(ctx context.Context) (availability []string, err error)
}
