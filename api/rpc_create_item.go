package api

import (
	"context"
	"log"
	"net/http"

	"github.com/lib/pq"
	"github.com/machearn/galaxy_service/api_error"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
)

func (server *Server) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	arg := db.CreateItemParams{
		Name:     req.Name,
		Quantity: req.Quantity,
		Price:    req.Price,
	}

	item, err := server.store.CreateItem(ctx, arg)
	if err != nil {
		pqErr := err.(*pq.Error)
		var apiErr *api_error.APIError
		if pqErr.Code[:2] == "23" {
			apiErr = api_error.NewAPIError(http.StatusBadRequest, "invalid input")
		} else {
			apiErr = api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		}
		log.Print("cannot create item: ", err)
		return nil, apiErr
	}

	res := pb.CreateItemResponse{
		Item: &pb.Item{
			ID:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		},
	}

	return &res, nil
}
