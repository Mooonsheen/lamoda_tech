package postgresql

import (
	"context"
	"fmt"

	"github.com/Mooonsheen/lamoda_tech/app/internal/models"
	"github.com/Mooonsheen/lamoda_tech/app/internal/storage"
	"github.com/Mooonsheen/lamoda_tech/app/internal/storage/interfaces"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type database struct {
	client storage.StorageClient
	pool   *pgxpool.Pool
}

func NewDatadase(client storage.StorageClient, pool *pgxpool.Pool) interfaces.Storage {
	return &database{
		client: client,
		pool:   pool,
	}
}

func (db *database) CreateReservation(ctx context.Context, msg models.ReserveRequestMessage, uuid string) (reservedItems models.ResponseMessageAfterCreate, err error) {
	q := `SELECT create_reservation_one_store($1, $2, $3, $4, $5)`

	var response models.ResponseMessageAfterCreate
	response.AppliedItems = make([]*models.ItemsMessage, len(msg.Items))

	conn, err := db.pool.Acquire(context.TODO())
	if err != nil {
		return response, err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(context.TODO(), pgx.TxOptions{})
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	var isReserved, countFaildReservations int
	for _, item := range msg.Items {
		tx.QueryRow(ctx, q, uuid, msg.ClientId, msg.Stores[0].Id, item.Id, item.Count).Scan(&isReserved)
		if isReserved == 0 {
			countFaildReservations++
			continue
		}
		response.AppliedItems = append(response.AppliedItems, item)
	}
	if len(msg.Items) == countFaildReservations {
		return response, fmt.Errorf("can't reservate all these items")
	}

	return response, nil
}

func (db *database) СonfirmReservation(ctx context.Context, uuid string) (status string, err error) {
	return "q", nil
}

func (db *database) DeleteReservation(ctx context.Context, uuid string) (status string, err error) {
	return "q", nil
}

func (db *database) GetStoreRemains(ctx context.Context, store_id string) (remains map[string]int, err error) {
	return nil, nil
}

func (db *database) GetStoresAvailability(ctx context.Context) (availability []string, err error) {
	q := `SELECT id
			FROM stores
			WHERE availability = false
			RETURNING id`

	rows, err := db.client.Query(ctx, q)
	if pgErr, ok := err.(*pgconn.PgError); ok {
		newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		logrus.Errorf("can't exec GetStoresAvailability database method: %s", newErr)
		ctx.Err()
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var store_id string
		err = rows.Scan(&store_id)
		if err != nil {
			logrus.Error(err)
		}
		availability = append(availability, store_id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return availability, nil
}

// func (db *database) CreateReserveSingleStore(ctx context.Context, msg models.Message) (error) {
// 	q := `BEGIN;

// 		COMMIT;`
// 	var productInfo model.ProductInfo
// 	err := db.client.QueryRow(ctx, q, id).Scan(&productInfo.ProductId, &productInfo.ProductName, &productInfo.ProductDescription)
// if pgErr, ok := err.(*pgconn.PgError); ok {
// 	newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
// 	logrus.Errorf("can't get table in ProductCRUD: %s", newErr)
// 	ctx.Err()
// } else if err != nil {
// 	return productInfo, err
// }

// 	return productInfo, nil
// }

// func (db *database) CreateReserveWithSeveralStores(ctx context.Context, msg models.Message) (model.ProductInfo, error) {

// }
