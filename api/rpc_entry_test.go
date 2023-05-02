package api

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/machearn/galaxy_service/db/mock"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"github.com/machearn/galaxy_service/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateRandomEntry() db.Entry {
	return db.Entry{
		ID:        int32(util.GetRandomInt(1, 20)),
		MemberID:  int32(util.GetRandomInt(1, 20)),
		ItemID:    int32(util.GetRandomInt(1, 20)),
		Quantity:  int32(util.GetRandomInt(1, 100)),
		Total:     int32(util.GetRandomInt(1, 1000)),
		CreatedAt: time.Now().UTC().Truncate(time.Second),
	}
}

func TestCreateEntryAPI(t *testing.T) {
	entry := CreateRandomEntry()

	req := pb.CreateEntryRequest{
		UserId:   entry.MemberID,
		ItemId:   entry.ItemID,
		Quantity: entry.Quantity,
		Total:    entry.Total,
	}

	arg := db.CreateEntryParams{
		MemberID:  req.GetUserId(),
		ItemID:    req.GetItemId(),
		Quantity:  req.GetQuantity(),
		Total:     req.GetTotal(),
		CreatedAt: entry.CreatedAt,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().CreateEntry(gomock.Any(), gomock.Eq(arg)).Times(1).Return(entry, nil)

	listener := startTestServer(t, store, nil)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	res, err := client.CreateEntry(context.Background(), &req)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.NotZero(t, res.GetEntry().GetID())
	require.Equal(t, entry.MemberID, res.GetEntry().GetUserId())
	require.Equal(t, entry.ItemID, res.GetEntry().GetItemId())
	require.Equal(t, entry.Quantity, res.GetEntry().GetQuantity())
	require.Equal(t, entry.Total, res.GetEntry().GetTotal())
	require.Equal(t, entry.CreatedAt, res.GetEntry().GetCreatedAt().AsTime())
}

func TestListEntriesAPI(t *testing.T) {
	n := 5
	entries := make([]db.Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = CreateRandomEntry()
	}

	req := pb.ListEntriesRequest{
		Limit:  5,
		Offset: 0,
	}

	arg := db.ListEntriesParams{
		Limit:  5,
		Offset: 0,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().ListEntries(gomock.Any(), gomock.Eq(arg)).Times(1).Return(entries, nil)

	listener := startTestServer(t, store, nil)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	res, err := client.ListEntries(context.Background(), &req)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Len(t, res.GetEntries(), 5)

	for i, entry := range res.GetEntries() {
		require.NotZero(t, entry.GetID())
		require.Equal(t, entry.GetUserId(), entries[i].MemberID)
		require.Equal(t, entry.GetItemId(), entries[i].ItemID)
		require.Equal(t, entry.GetQuantity(), entries[i].Quantity)
		require.Equal(t, entry.GetTotal(), entries[i].Total)
		require.Equal(t, entry.GetCreatedAt().AsTime(), entries[i].CreatedAt)
	}
}

func TestListEntriesByUserAPI(t *testing.T) {
	n := 5
	entries := make([]db.Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = CreateRandomEntry()
		entries[i].MemberID = 1
	}

	req := pb.ListEntriesByUserRequest{
		UserId: entries[0].MemberID,
		Limit:  5,
		Offset: 0,
	}

	arg := db.ListEntriesByMemberParams{
		MemberID: entries[0].MemberID,
		Limit:    5,
		Offset:   0,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().ListEntriesByMember(gomock.Any(), gomock.Eq(arg)).Times(1).Return(entries, nil)

	listener := startTestServer(t, store, nil)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)
	res, err := client.ListEntriesByUser(context.Background(), &req)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Len(t, res.GetEntries(), 5)
	for i, entry := range res.GetEntries() {
		require.NotZero(t, entry.GetID())
		require.Equal(t, entry.GetUserId(), entries[i].MemberID)
		require.Equal(t, entry.GetItemId(), entries[i].ItemID)
		require.Equal(t, entry.GetQuantity(), entries[i].Quantity)
		require.Equal(t, entry.GetTotal(), entries[i].Total)
		require.Equal(t, entry.GetCreatedAt().AsTime(), entries[i].CreatedAt)
	}
}

func TestListEntriesByItemAPI(t *testing.T) {
	n := 5
	entries := make([]db.Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = CreateRandomEntry()
		entries[i].ItemID = 1
	}

	req := pb.ListEntriesByItemRequest{
		ItemId: entries[0].ItemID,
		Limit:  5,
		Offset: 0,
	}

	arg := db.ListEntriesByItemParams{
		ItemID: entries[0].ItemID,
		Limit:  5,
		Offset: 0,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().ListEntriesByItem(gomock.Any(), gomock.Eq(arg)).Times(1).Return(entries, nil)

	listener := startTestServer(t, store, nil)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)
	res, err := client.ListEntriesByItem(context.Background(), &req)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Len(t, res.GetEntries(), 5)
	for i, entry := range res.GetEntries() {
		require.NotZero(t, entry.GetID())
		require.Equal(t, entry.GetUserId(), entries[i].MemberID)
		require.Equal(t, entry.GetItemId(), entries[i].ItemID)
		require.Equal(t, entry.GetQuantity(), entries[i].Quantity)
		require.Equal(t, entry.GetTotal(), entries[i].Total)
		require.Equal(t, entry.GetCreatedAt().AsTime(), entries[i].CreatedAt)
	}
}

func TestDeleteEntryAPI(t *testing.T) {
	entry := CreateRandomEntry()

	req := pb.DeleteEntryRequest{
		Id: entry.ID,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().DeleteEntry(gomock.Any(), gomock.Eq(entry.ID)).Times(1).Return(nil)

	listener := startTestServer(t, store, nil)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)
	res, err := client.DeleteEntry(context.Background(), &req)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.True(t, res.GetSuccess())
}
