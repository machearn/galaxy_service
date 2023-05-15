package api

import (
	"context"
	"database/sql"
	"log"

	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.GetItemResponse, error) {
	item, err := server.store.GetItem(ctx, req.GetId())
	if err != nil {
		log.Print("cannot get item: ", err)
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "item not found: %v", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	res := pb.GetItemResponse{
		Item: &pb.Item{
			ID:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		},
	}

	return &res, nil
}

func (server *Server) ListItems(ctx context.Context, req *pb.ListItemsRequest) (*pb.ListItemsResponse, error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	arg := db.ListItemsParams{
		Offset: offset,
		Limit:  limit,
	}

	rows, err := server.store.ListItems(ctx, arg)
	if err != nil {
		log.Print("cannot get item: ", err)
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "item not found: %v", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	var items []*pb.Item
	for _, item := range rows {
		items = append(items, &pb.Item{
			ID:       item.ID,
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}

	rep := pb.ListItemsResponse{
		Items: items,
	}

	return &rep, nil
}
