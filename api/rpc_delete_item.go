package api

import (
	"context"
	"log"

	"github.com/lib/pq"
	"github.com/machearn/galaxy_service/pb"
)

func (server *Server) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	err := server.store.DeleteItem(ctx, req.GetId())
	if err != nil {
		pqErr := err.(*pq.Error)
		log.Print("cannot delete item: ", pqErr)
		return &pb.DeleteItemResponse{
			Success: false,
		}, pqErr
	}

	rep := pb.DeleteItemResponse{
		Success: true,
	}

	return &rep, nil
}
