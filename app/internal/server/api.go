package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Mooonsheen/lamoda_tech/app/internal/models"
	"github.com/Mooonsheen/lamoda_tech/app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/jsonapi"
)

func (s *Server) handleReservation(ctx *gin.Context) {
	switch ctx.Request.Method {
	case http.MethodPost:
		s.handlePostReservation(ctx)
	case http.MethodPatch:
		s.handlePatchReservation(ctx)
	case http.MethodDelete:
		s.handleDeleteReservation(ctx)
	default:
		http.Error(ctx.Writer, "not found", http.StatusNotFound)
		return
	}
}

func (s *Server) handlePostReservation(ctx *gin.Context) {
	msg := new(models.RequestReserveMessage)

	ctx.Writer.Header().Set("Content-Type", jsonapi.MediaType)
	ctx.Writer.WriteHeader(http.StatusCreated)

	response := new(models.ResponseReserveCreateMessage)

	if err := jsonapi.UnmarshalPayload(ctx.Request.Body, msg); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		jsonapi.MarshalPayload(ctx.Writer, response)
		return
	}
	err := validateRequestPostMessage(msg)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		jsonapi.MarshalPayload(ctx.Writer, response)
		return
	}

	err = PgClient.GetStoresAvailability(ctx, *msg)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		jsonapi.MarshalPayload(ctx.Writer, response)
		return
	}

	uuid := utils.GenerateUUID()

	if len(msg.Stores) > 1 {
		response, err = PgClient.CreateReservationManyStore(ctx, *msg, uuid)
	} else {
		response, err = PgClient.CreateReservation(ctx, *msg, uuid)
	}
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		jsonapi.MarshalPayload(ctx.Writer, response)
		return
	}

	if err := jsonapi.MarshalPayload(ctx.Writer, response); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handlePatchReservation(ctx *gin.Context) {
	msg := new(models.RequestReserveMessage)

	ctx.Writer.Header().Set("Content-Type", jsonapi.MediaType)
	ctx.Writer.WriteHeader(http.StatusBadRequest)

	response := new(models.ResponseReserveChangeMessage)
	response.Status = "0"

	if err := jsonapi.UnmarshalPayload(ctx.Request.Body, msg); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		jsonapi.MarshalPayload(ctx.Writer, response)
		return
	}
	err := validateRequestPatchDeleteMessage(*msg)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		jsonapi.MarshalPayload(ctx.Writer, response)
		return
	}

	response, err = PgClient.ApplyReservation(ctx, *msg)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		jsonapi.MarshalPayload(ctx.Writer, response)
		return
	}

	ctx.Writer.Header().Set("Content-Type", jsonapi.MediaType)
	ctx.Writer.WriteHeader(http.StatusCreated)

	if err := jsonapi.MarshalPayload(ctx.Writer, response); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleDeleteReservation(ctx *gin.Context) {
	msg := new(models.RequestReserveMessage)

	ctx.Writer.Header().Set("Content-Type", jsonapi.MediaType)
	ctx.Writer.WriteHeader(http.StatusBadRequest)

	response := new(models.ResponseReserveChangeMessage)
	response.Status = "0"

	if err := jsonapi.UnmarshalPayload(ctx.Request.Body, msg); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		jsonapi.MarshalPayload(ctx.Writer, response)
		return
	}

	err := validateRequestPatchDeleteMessage(*msg)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		jsonapi.MarshalPayload(ctx.Writer, response)
		return
	}

	response, err = PgClient.DeleteReservation(ctx, *msg)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		jsonapi.MarshalPayload(ctx.Writer, response)
		return
	}

	ctx.Writer.Header().Set("Content-Type", jsonapi.MediaType)
	ctx.Writer.WriteHeader(http.StatusCreated)

	if err := jsonapi.MarshalPayload(ctx.Writer, response); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleGetStoreRemains(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", jsonapi.MediaType)
	ctx.Writer.WriteHeader(http.StatusBadRequest)

	response := new(models.StoreMessage)
	response.Id = ctx.Query("id")
	sortType := ctx.Query("sort")
	pages := ctx.Query("pages")

	currentSortType, currentPages, err := validateGetStoreRemains(response.Id, sortType, pages)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		jsonapi.MarshalPayload(ctx.Writer, response)
		return
	}

	response, err = PgClient.GetStoreRemains(ctx, response.Id, currentSortType, currentPages)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		jsonapi.MarshalPayload(ctx.Writer, response)
		return
	}

	ctx.Writer.Header().Set("Content-Type", jsonapi.MediaType)
	ctx.Writer.WriteHeader(http.StatusCreated)

	if err := jsonapi.MarshalPayload(ctx.Writer, response); err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func validateRequestPostMessage(msg *models.RequestReserveMessage) error {
	if len(msg.Items) != 0 {
		log.Printf("len(msg.Items): %d", len(msg.Items))
		for i := 0; i < len(msg.Items); i++ {
			if msg.Items[i].Amount <= 0 {
				return fmt.Errorf("you should order a positive count of items")
			}
		}
	} else {
		return fmt.Errorf("you should order at least one item")
	}
	if len(msg.Stores) == 0 {
		log.Printf("len(msg.Stores): %d", len(msg.Stores))
		availableStores, err := PgClient.GetAllAvailableStores(context.TODO())
		if err != nil {
			return fmt.Errorf("can't get all available stores")
		}
		msg.Stores = availableStores
	}
	log.Printf("msg: %v", *msg)
	return nil
}

func validateRequestPatchDeleteMessage(msg models.RequestReserveMessage) error {
	if len(msg.Uiud) < 1 {
		return fmt.Errorf("you should choose the uuid of your order")
	}
	return nil
}

func validateGetStoreRemains(storeId, sortType, pages string) (currentSortType string, currentPages int, err error) {
	if len(storeId) == 0 {
		return currentSortType, 0, fmt.Errorf("you should choose one store")
	}
	if len(pages) == 0 {
		currentPages = 1
	} else {
		p, err := strconv.Atoi(pages)
		if err != nil {
			return currentSortType, 0, fmt.Errorf("invalid pages param")
		}
		if p < 0 {
			return currentSortType, 0, fmt.Errorf("pages param should be a positive number")
		}
		currentPages = p
	}
	if len(sortType) == 0 {
		currentSortType = "desc"
		return currentSortType, currentPages, nil
	} else if sortType == "asc" || sortType == "desc" {
		return sortType, currentPages, nil
	}
	return currentSortType, 0, fmt.Errorf("invalid sort param, it should be literally asc or desc")
}
