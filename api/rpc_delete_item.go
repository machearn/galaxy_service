package api

import (
	"context"
	"log"

	"github.com/machearn/galaxy_service/pb"
)

func (server *Server) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.DeleteItemResponse, error) {
	err := server.store.DeleteItem(ctx, req.GetId())
	if err != nil {
		log.Print("cannot delete item: ", err)
		return &pb.DeleteItemResponse{
			Success: false,
		}, err
	}

	rep := pb.DeleteItemResponse{
		Success: true,
	}

	return &rep, nil
}
