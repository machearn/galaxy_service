package api

import (
	"context"
	"database/sql"
	"log"

	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.Empty, error) {
	err := server.store.DeleteItem(ctx, req.GetId())
	if err != nil {
		log.Print("cannot delete item: ", err)
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.InvalidArgument, "item not found: %v", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	return &pb.Empty{}, nil
}
