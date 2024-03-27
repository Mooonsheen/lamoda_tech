package server

import (
	"log"
	"net/http"

	"github.com/Mooonsheen/lamoda_tech/app/internal/models"
	"github.com/Mooonsheen/lamoda_tech/app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
)

/*

кэш для зарезервированных товаров (map[client_name]map[item]reserved)

1) Reserve items in one store
In: (id/count map[int]int)
Out: 200 or err ()

2) Unreserve items in one store (реализовать проверку по кэшу ) (сделать хранимку в бд на апдейт после отмены)
In: (id/count map[int]int)
Out: 200 or err ()

3) Check the leftovers of all items in one store (сделать пагинацию на выдачу)
In: (store_id int)
Out: 200 or err (если некорректный id склада)


*/

func (s *Server) handleReservation(ctx *gin.Context) {
	// if /*r.Header.Get(headerAccept)*/ ctx.Request.Header.Values("Conte")[] != jsonapi.MediaType {
	// 	http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
	// }

	switch ctx.Request.Method {
	case http.MethodPost:
		s.handlePostReservation(ctx)
	// case http.MethodPatch:
	// 	s.handlePatchReservation(ctx)
	case http.MethodDelete:
		s.handleDeleteReservation(ctx)
	default:
		http.Error(ctx.Writer, "Not Found", http.StatusNotFound)
		return
	}
}

func (s *Server) handlePostReservation(ctx *gin.Context) {
	msg := new(models.ReserveRequestMessage)

	if err := jsonapi.UnmarshalPayload(ctx.Request.Body, msg); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// pgClient := postgresql.NewDatadase(Pool)
	uuid := utils.GenerateUUID()
	reservedItems, err := PgClient.CreateReservation(ctx, *msg, uuid)
	if err != nil {
		log.Printf("can't reserve in handler handlePostReservation")
	}
	reservedItems.Id = uuid

	//InternalCache.AddReserve(*msg)
	// cache := InternalCache.GetReservation()
	// log.Print(cache)

	ctx.Writer.Header().Set("Content-Type", jsonapi.MediaType)
	ctx.Writer.WriteHeader(http.StatusCreated)

	if err := jsonapi.MarshalPayload(ctx.Writer, reservedItems); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleDeleteReservation(ctx *gin.Context) {

}

// func (s *Server) handleGetStoreRemains(ctx *gin.Context) {
// 	// validate content-type
// 	if err := jsonapi.UnmarshalPayload(ctx.Request.Body, models.Message); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// поход в бд за данными
// 	pqStorage := postgresql.NewDatadase(Pool)

// 	w.Header().Set("Content-Type", jsonapi.MediaType)
// 	w.WriteHeader(http.StatusCreated)

// 	if err := jsonapi.MarshalPayload(w, blog); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// func validateStore(id int) (statusCode, error) {

// }
