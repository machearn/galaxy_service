package api

import (
	"context"
	"database/sql"
	"log"

	"github.com/lib/pq"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		log.Print("cannot update item: ", err)
		pqErr := err.(*pq.Error)
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "item not found: %v", err.Error())
		}
		if pqErr.Code[:2] == "23" {
			return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
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
