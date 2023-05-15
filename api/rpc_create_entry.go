package api

import (
	"context"
	"log"
	"time"

	"github.com/lib/pq"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateEntry(ctx context.Context, req *pb.CreateEntryRequest) (*pb.CreateEntryResponse, error) {
	arg := db.CreateEntryParams{
		MemberID:  req.GetUserId(),
		ItemID:    req.GetItemId(),
		Quantity:  req.GetQuantity(),
		Total:     req.GetTotal(),
		CreatedAt: time.Now().UTC().Truncate(time.Second),
	}

	entry, err := server.store.CreateEntry(ctx, arg)
	if err != nil {
		log.Print("cannot create entry: ", err)
		pqErr := err.(*pq.Error)
		if pqErr.Code[:2] == "23" {
			return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	res := pb.CreateEntryResponse{
		Entry: &pb.Entry{
			ID:        entry.ID,
			UserId:    entry.MemberID,
			ItemId:    entry.ItemID,
			Quantity:  entry.Quantity,
			Total:     entry.Total,
			CreatedAt: timestamppb.New(entry.CreatedAt),
		},
	}

	return &res, nil
}
