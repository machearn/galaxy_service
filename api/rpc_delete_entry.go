package api

import (
	"context"
	"log"

	"github.com/lib/pq"
	"github.com/machearn/galaxy_service/pb"
)

func (server *Server) DeleteEntry(ctx context.Context, req *pb.DeleteEntryRequest) (*pb.DeleteEntryResponse, error) {
	err := server.store.DeleteEntry(ctx, req.GetId())
	if err != nil {
		pqErr := err.(*pq.Error)
		log.Print("cannot delete entry: ", pqErr)
		return &pb.DeleteEntryResponse{
			Success: false,
		}, pqErr
	}

	rep := pb.DeleteEntryResponse{
		Success: true,
	}

	return &rep, nil
}
