package api

import (
	"context"
	"log"

	"github.com/lib/pq"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.CreateItemResponse, error) {
	arg := db.CreateItemParams{
		Name:     req.Name,
		Quantity: req.Quantity,
		Price:    req.Price,
	}

	item, err := server.store.CreateItem(ctx, arg)
	if err != nil {
		log.Print("cannot create item: ", err)
		pqErr := err.(*pq.Error)
		if pqErr.Code[:2] == "23" {
			return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
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
