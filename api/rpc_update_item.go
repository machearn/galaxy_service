package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/lib/pq"
	"github.com/machearn/galaxy_service/api_error"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
)

func (server *Server) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.UpdateItemResponse, error) {
	arg := db.UpdateItemParams{
		ID: req.GetId(),
	}

	if req.Name != nil {
		arg.Name = sql.NullString{
			String: req.GetName(),
			Valid:  true,
		}
	}

	if req.Quantity != nil {
		arg.Quantity = sql.NullInt32{
			Int32: req.GetQuantity(),
			Valid: true,
		}
	}

	if req.Price != nil {
		arg.Price = sql.NullInt32{
			Int32: req.GetPrice(),
			Valid: true,
		}
	}

	item, err := server.store.UpdateItem(ctx, arg)
	if err != nil {
		pqErr := err.(*pq.Error)
		var apiErr *api_error.APIError
		if err == sql.ErrNoRows {
			apiErr = api_error.NewAPIError(http.StatusBadRequest, "item not found")
		} else if pqErr.Code[:2] == "23" {
			apiErr = api_error.NewAPIError(http.StatusBadRequest, "invalid input")
		} else {
			apiErr = api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		}
		log.Print("cannot update item: ", err)
		return nil, apiErr
	}

	res := pb.UpdateItemResponse{
		Item: &pb.Item{
			ID:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		},
	}

	return &res, nil
}
