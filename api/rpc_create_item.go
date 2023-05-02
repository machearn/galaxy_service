package api

import (
	"context"
	"log"

	"github.com/lib/pq"
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
		log.Print("cannot create item: ", pqErr)
		return nil, pqErr
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
