package api

import (
	"context"
	"log"

	"github.com/lib/pq"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) GetEntry(ctx context.Context, req *pb.GetEntryRequest) (*pb.GetEntryResponse, error) {
	ID := req.GetId()

	entry, err := server.store.GetEntry(ctx, ID)
	if err != nil {
		pqErr := err.(*pq.Error)
		log.Print("failed to get entry: ", pqErr)
		return nil, pqErr
	}

	res := pb.GetEntryResponse{
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

func (server *Server) ListEntries(ctx context.Context, req *pb.ListEntriesRequest) (*pb.ListEntriesResponse, error) {
	arg := db.ListEntriesParams{
		Offset: req.GetOffset(),
		Limit:  req.GetLimit(),
	}

	rows, err := server.store.ListEntries(ctx, arg)
	if err != nil {
		pqErr := err.(*pq.Error)
		log.Print("failed to list entries: ", pqErr)
		return nil, pqErr
	}

	var entries []*pb.Entry
	for _, entry := range rows {
		entries = append(entries, &pb.Entry{
			ID:        entry.ID,
			UserId:    entry.MemberID,
			ItemId:    entry.ItemID,
			Quantity:  entry.Quantity,
			Total:     entry.Total,
			CreatedAt: timestamppb.New(entry.CreatedAt),
		})
	}

	res := pb.ListEntriesResponse{
		Entries: entries,
	}

	return &res, nil
}

func (server *Server) ListEntriesByUser(ctx context.Context, req *pb.ListEntriesByUserRequest) (*pb.ListEntriesResponse, error) {
	arg := db.ListEntriesByMemberParams{
		MemberID: req.GetUserId(),
		Offset:   req.GetOffset(),
		Limit:    req.GetLimit(),
	}

	rows, err := server.store.ListEntriesByMember(ctx, arg)
	if err != nil {
		pqErr := err.(*pq.Error)
		log.Print("failed to list entries: ", pqErr)
		return nil, pqErr
	}

	var entries []*pb.Entry
	for _, entry := range rows {
		entries = append(entries, &pb.Entry{
			ID:        entry.ID,
			UserId:    entry.MemberID,
			ItemId:    entry.ItemID,
			Quantity:  entry.Quantity,
			Total:     entry.Total,
			CreatedAt: timestamppb.New(entry.CreatedAt),
		})
	}

	res := pb.ListEntriesResponse{
		Entries: entries,
	}

	return &res, nil
}

func (server *Server) ListEntriesByItem(ctx context.Context, req *pb.ListEntriesByItemRequest) (*pb.ListEntriesResponse, error) {
	arg := db.ListEntriesByItemParams{
		ItemID: req.GetItemId(),
		Offset: req.GetOffset(),
		Limit:  req.GetLimit(),
	}

	rows, err := server.store.ListEntriesByItem(ctx, arg)
	if err != nil {
		pqErr := err.(*pq.Error)
		log.Print("failed to list entries: ", pqErr)
		return nil, pqErr
	}

	var entries []*pb.Entry
	for _, entry := range rows {
		entries = append(entries, &pb.Entry{
			ID:        entry.ID,
			UserId:    entry.MemberID,
			ItemId:    entry.ItemID,
			Quantity:  entry.Quantity,
			Total:     entry.Total,
			CreatedAt: timestamppb.New(entry.CreatedAt),
		})
	}

	res := pb.ListEntriesResponse{
		Entries: entries,
	}

	return &res, nil
}
