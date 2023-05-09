package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/lib/pq"
	"github.com/machearn/galaxy_service/api_error"
	"github.com/machearn/galaxy_service/pb"
)

func (server *Server) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.Empty, error) {
	err := server.store.DeleteItem(ctx, req.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			apiErr := api_error.NewAPIError(http.StatusBadGateway, "item not found")
			return nil, apiErr
		}
		pqErr := err.(*pq.Error)
		log.Print("cannot delete item: ", pqErr)
		return nil, err
	}

	return &pb.Empty{}, nil
}
