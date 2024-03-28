package interfaces

import (
	"context"

	"github.com/Mooonsheen/lamoda_tech/app/internal/models"
)

type Storage interface {
	CreateReservation(ctx context.Context, msg models.RequestReserveMessage, uuid string) (*models.ResponseReserveCreateMessage, error)
	CreateReservationManyStore(ctx context.Context, msg models.RequestReserveMessage, uuid string) (*models.ResponseReserveCreateMessage, error)
	ApplyReservation(ctx context.Context, msg models.RequestReserveMessage) (*models.ResponseReserveChangeMessage, error)
	DeleteReservation(ctx context.Context, msg models.RequestReserveMessage) (*models.ResponseReserveChangeMessage, error)
	GetStoreRemains(ctx context.Context, store_id, currentSortType string, currentString int) (*models.StoreMessage, error)
	GetStoresAvailability(ctx context.Context, msg models.RequestReserveMessage) error
	GetAllAvailableStores(ctx context.Context) ([]*models.StoreMessage, error)
}
