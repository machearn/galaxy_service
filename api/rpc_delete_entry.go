package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/machearn/galaxy_service/api_error"
	"github.com/machearn/galaxy_service/pb"
)

func (server *Server) DeleteEntry(ctx context.Context, req *pb.DeleteEntryRequest) (*pb.Empty, error) {
	err := server.store.DeleteEntry(ctx, req.GetId())
	if err != nil {
		var apiErr *api_error.APIError
		if err == sql.ErrNoRows {
			apiErr = api_error.NewAPIError(http.StatusBadGateway, "entry not found")
		} else {
			apiErr = api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		}
		log.Print("cannot delete entry: ", err)
		return nil, apiErr
	}

	return &pb.Empty{}, nil
}
