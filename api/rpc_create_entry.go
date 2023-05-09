package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
	"github.com/machearn/galaxy_service/api_error"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
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
		pqErr := err.(*pq.Error)
		var apiErr *api_error.APIError
		if pqErr.Code[:2] == "23" {
			apiErr = api_error.NewAPIError(http.StatusBadRequest, "invalid input")
		} else {
			apiErr = api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		}
		log.Print("cannot create entry: ", err)
		return nil, apiErr
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
