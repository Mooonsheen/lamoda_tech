package postgresql

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Mooonsheen/lamoda_tech/app/internal/models"
	"github.com/Mooonsheen/lamoda_tech/app/internal/storage/interfaces"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type database struct {
	pool *pgxpool.Pool
}

func NewDatadase(pool *pgxpool.Pool) interfaces.Storage {
	return &database{
		pool: pool,
	}
}

func (db *database) CreateReservation(ctx context.Context, msg models.RequestReserveMessage, uuid string) (*models.ResponseReserveCreateMessage, error) {
	q := `SELECT create_reservation_one_store($1, $2, $3, $4, $5)`

	response := new(models.ResponseReserveCreateMessage)
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
	notReservedItems := make(map[string][]int)
	for _, item := range msg.Items {
		tx.QueryRow(ctx, q, uuid, msg.ClientId, msg.Stores[0].Id, item.Id, item.Amount).Scan(&isReserved)
		if isReserved == 0 {
			countFaildReservations++
			subQuery := `SELECT total_amount FROM item_total_amount WHERE item_id = $1`
			var itemTotalAmount, itemLocalStoreAmount int
			tx.QueryRow(ctx, subQuery, item.Id).Scan(&itemTotalAmount)
			subQuery = `SELECT amount FROM store_item WHERE store_id = $1 AND item_id = $2`
			tx.QueryRow(ctx, subQuery, msg.Stores[0].Id, item.Id).Scan(&itemLocalStoreAmount)
			notReservedItems[item.Id] = []int{item.Amount, itemLocalStoreAmount, itemTotalAmount}
			continue
		}
		response.AppliedItems = append(response.AppliedItems, &models.Items{
			Id:         item.Id,
			ReservedIn: msg.Stores[0].Id,
			Amount:     item.Amount,
		})
	}
	if countFaildReservations != 0 {
		tx.Rollback(context.TODO())
		var explain string
		for item_id, amounts := range notReservedItems {
			explain += ("id: " + item_id + ", required: " + strconv.Itoa(amounts[0]) + ", store amount: " + strconv.Itoa(amounts[1]) + ", total amount: " + strconv.Itoa(amounts[2]) + "; ")
		}
		return &models.ResponseReserveCreateMessage{}, fmt.Errorf("can't reservate all these items, in this explane you can see problem items.\n%s", explain)
	}
	response.Id = uuid

	return response, nil
}

func (db *database) CreateReservationManyStore(ctx context.Context, msg models.RequestReserveMessage, uuid string) (*models.ResponseReserveCreateMessage, error) {
	q := `SELECT create_reservation_many_stores($1, $2, $3, $4, $5)`

	response := new(models.ResponseReserveCreateMessage)
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

	storesForReservate := make([]string, 0, len(msg.Stores))
	for _, store := range msg.Stores {
		storesForReservate = append(storesForReservate, store.Id)
	}

	var isReserved, countFaildReservations int
	notReservedItems := make(map[string][]int)
	for _, item := range msg.Items {
		tx.QueryRow(ctx, q, uuid, msg.ClientId, storesForReservate, item.Id, item.Amount).Scan(&isReserved)
		if isReserved != 1 {
			countFaildReservations++
			subQuery := `SELECT total_amount FROM item_total_amount WHERE item_id = $1`
			var itemTotalAmount int
			tx.QueryRow(ctx, subQuery, item.Id).Scan(&itemTotalAmount)
			notReservedItems[item.Id] = []int{itemTotalAmount, item.Amount}
			continue
		}
		response.AppliedItems = append(response.AppliedItems, &models.Items{
			Id:         item.Id,
			ReservedIn: msg.Stores[0].Id,
			Amount:     item.Amount,
		})
	}
	if countFaildReservations != 0 {
		tx.Rollback(context.TODO())
		var explain string
		for item_id, amounts := range notReservedItems {
			explain += ("id: " + item_id + ", total_amount: " + strconv.Itoa(amounts[0]) + ", required: " + strconv.Itoa(amounts[1]) + "; ")
		}
		return &models.ResponseReserveCreateMessage{}, fmt.Errorf("can't reservate all these items, in this explane you can see problem items.\n%s", explain)
	}

	response.Id = uuid

	return response, nil
}

func (db *database) ApplyReservation(ctx context.Context, msg models.RequestReserveMessage) (*models.ResponseReserveChangeMessage, error) {
	q := `SELECT apply_reservation_one_store($1, $2)`

	response := new(models.ResponseReserveChangeMessage)
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

	var isApplied int
	tx.QueryRow(ctx, q, msg.ClientId, msg.Uiud).Scan(&isApplied)
	if isApplied == 0 {
		response = &models.ResponseReserveChangeMessage{
			Id:     msg.Id,
			Status: "0",
		}
		return response, fmt.Errorf("you have not this order or it's already have deleted/applied")
	} else {
		response = &models.ResponseReserveChangeMessage{
			Id:     msg.Id,
			Status: "1",
		}
	}

	return response, nil
}

func (db *database) DeleteReservation(ctx context.Context, msg models.RequestReserveMessage) (*models.ResponseReserveChangeMessage, error) {
	q := `SELECT delete_reservation_one_store($1, $2)`

	response := new(models.ResponseReserveChangeMessage)
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

	var isDeleted int
	tx.QueryRow(ctx, q, msg.ClientId, msg.Uiud).Scan(&isDeleted)
	if isDeleted == 0 {
		response = &models.ResponseReserveChangeMessage{
			Id:     msg.Id,
			Status: "0",
		}
		return response, fmt.Errorf("you have not this order or it's already have deleted/applied")
	} else {
		response = &models.ResponseReserveChangeMessage{
			Id:     msg.Id,
			Status: "1",
		}
	}

	return response, nil
}

func (db *database) GetStoreRemains(ctx context.Context, storeId, currentSortType string, currentString int) (*models.StoreMessage, error) {
	q := `SELECT * FROM get_store_remains($1, $2, $3)`

	response := new(models.StoreMessage)
	response.Id = storeId
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

	rows, err := tx.Query(ctx, q, storeId, currentSortType, currentString)
	if err != nil {
		return response, fmt.Errorf("store with this id does not exist: %s", storeId)
	}
	defer rows.Close()

	for rows.Next() {
		var item_id string
		var amount int
		err = rows.Scan(&item_id, &amount)
		if err != nil {
			return response, fmt.Errorf("store with this id does not exist: %s", storeId)
		}
		response.Items = append(response.Items, &models.Items{
			Id:     item_id,
			Amount: amount,
		})
	}

	return response, nil
}

func (db *database) GetStoresAvailability(ctx context.Context, msg models.RequestReserveMessage) error {
	q := `SELECT check_store_availability($1)`

	conn, err := db.pool.Acquire(context.TODO())
	if err != nil {
		return err
	}
	defer conn.Release()

	var isAvailable bool
	notAvailableStores := make([]string, 0, len(msg.Stores))
	for _, store := range msg.Stores {
		conn.QueryRow(ctx, q, store.Id).Scan(&isAvailable)
		if !isAvailable {
			notAvailableStores = append(notAvailableStores, store.Id)
		}
	}
	if len(notAvailableStores) != 0 {
		var invalidStores string
		for _, store_id := range notAvailableStores {
			invalidStores += (store_id + ", ")
		}
		return fmt.Errorf("stores with these ids are not avaliable: %s please, choose other stores", invalidStores)
	}
	return nil
}

func (db *database) GetAllAvailableStores(ctx context.Context) ([]*models.StoreMessage, error) {
	q := `SELECT get_all_available_stores()`

	conn, err := db.pool.Acquire(context.TODO())
	if err != nil {
		return []*models.StoreMessage{}, err
	}
	defer conn.Release()

	response := make([]*models.StoreMessage, 0)

	rows, err := conn.Query(ctx, q)
	if err != nil {
		return response, fmt.Errorf("can't get all available stores")
	}
	defer rows.Close()

	for rows.Next() {
		var store_id string
		var status bool
		err = rows.Scan(&store_id, &status)
		if err != nil {
			return response, fmt.Errorf("can't get all available stores")
		}
		if status {
			response = append(response, &models.StoreMessage{
				Id:          store_id,
				IsAvailable: status,
			})
		}

	}

	return response, nil
}
