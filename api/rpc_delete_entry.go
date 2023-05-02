package api

import (
	"context"
	"log"

	"github.com/machearn/galaxy_service/pb"
)

func (server *Server) DeleteEntry(ctx context.Context, req *pb.DeleteEntryRequest) (*pb.DeleteEntryResponse, error) {
	err := server.store.DeleteEntry(ctx, req.GetId())
	if err != nil {
		log.Print("cannot delete entry: ", err)
		return &pb.DeleteEntryResponse{
			Success: false,
		}, err
	}

	rep := pb.DeleteEntryResponse{
		Success: true,
	}

	return &rep, nil
}
