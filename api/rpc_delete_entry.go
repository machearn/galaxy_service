package api

import (
	"context"
	"database/sql"
	"log"

	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) DeleteEntry(ctx context.Context, req *pb.DeleteEntryRequest) (*pb.Empty, error) {
	err := server.store.DeleteEntry(ctx, req.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.InvalidArgument, "entry not found: %v", err.Error())
		}
		log.Print("cannot delete entry: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	return &pb.Empty{}, nil
}
