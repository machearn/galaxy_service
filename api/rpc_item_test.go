package api

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/machearn/galaxy_service/db/mock"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"github.com/machearn/galaxy_service/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func creatRandomItem() db.Item {
	item := db.Item{
		ID:       int32(util.GetRandomInt(1, 100)),
		Name:     util.GetRandomString(10),
		Quantity: int32(util.GetRandomInt(1, 100)),
		Price:    int32(util.GetRandomInt(1, 100)),
	}

	return item
}

func TestCreateItemAPI(t *testing.T) {
	item := creatRandomItem()

	req := pb.CreateItemRequest{
		Name:     item.Name,
		Quantity: item.Quantity,
		Price:    item.Price,
	}

	arg := db.CreateItemParams{
		Name:     req.Name,
		Quantity: req.Quantity,
		Price:    req.Price,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().CreateItem(gomock.Any(), gomock.Eq(arg)).Times(1).Return(item, nil)

	listener := startTestServer(t, store)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	res, err := client.CreateItem(context.Background(), &req)
	require.NoError(t, err)
	require.Equal(t, item.ID, res.GetItem().ID)
	require.Equal(t, item.Name, res.GetItem().Name)
	require.Equal(t, item.Quantity, res.GetItem().Quantity)
	require.Equal(t, item.Price, res.GetItem().Price)
}

func TestGetItemAPI(t *testing.T) {
	item := creatRandomItem()

	req := pb.GetItemRequest{
		Id: item.ID,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().GetItem(gomock.Any(), gomock.Eq(item.ID)).Times(1).Return(item, nil)

	listener := startTestServer(t, store)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	res, err := client.GetItem(context.Background(), &req)
	require.NoError(t, err)
	require.Equal(t, item.ID, res.GetItem().ID)
	require.Equal(t, item.Name, res.GetItem().Name)
	require.Equal(t, item.Quantity, res.GetItem().Quantity)
	require.Equal(t, item.Price, res.GetItem().Price)
}

func TestListItemsAPI(t *testing.T) {
	var items []db.Item

	for i := 0; i < 5; i++ {
		items = append(items, creatRandomItem())
	}

	req := pb.ListItemsRequest{
		Offset: 0,
		Limit:  5,
	}

	arg := db.ListItemsParams{
		Offset: req.GetOffset(),
		Limit:  req.GetLimit(),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().ListItems(gomock.Any(), gomock.Eq(arg)).Times(1).Return(items, nil)

	listener := startTestServer(t, store)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	res, err := client.ListItems(context.Background(), &req)
	require.NoError(t, err)
	require.Len(t, res.GetItems(), 5)

	for index, item := range items {
		require.Equal(t, items[index].ID, item.ID)
		require.Equal(t, items[index].Name, item.Name)
		require.Equal(t, items[index].Quantity, item.Quantity)
		require.Equal(t, items[index].Price, item.Price)
	}
}

func TestUpdateItemAPI(t *testing.T) {
	item := creatRandomItem()

	newItem := db.Item{
		ID:       item.ID,
		Name:     "New Name",
		Quantity: 100,
		Price:    100,
	}

	req := pb.UpdateItemRequest{
		Id:       item.ID,
		Name:     &newItem.Name,
		Quantity: &newItem.Quantity,
		Price:    &newItem.Price,
	}

	arg := db.UpdateItemParams{
		ID:       req.GetId(),
		Name:     sql.NullString{String: newItem.Name, Valid: true},
		Quantity: sql.NullInt32{Int32: newItem.Quantity, Valid: true},
		Price:    sql.NullInt32{Int32: newItem.Price, Valid: true},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().UpdateItem(gomock.Any(), gomock.Eq(arg)).Times(1).Return(newItem, nil)

	listener := startTestServer(t, store)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	res, err := client.UpdateItem(context.Background(), &req)
	require.NoError(t, err)
	require.Equal(t, newItem.ID, res.GetItem().ID)
	require.Equal(t, newItem.Name, res.GetItem().Name)
	require.Equal(t, newItem.Quantity, res.GetItem().Quantity)
	require.Equal(t, newItem.Price, res.GetItem().Price)
}

func TestDeleteItemAPI(t *testing.T) {
	item := creatRandomItem()

	req := pb.DeleteItemRequest{
		Id: item.ID,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().DeleteItem(gomock.Any(), gomock.Eq(item.ID)).Times(1).Return(nil)

	listener := startTestServer(t, store)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	res, err := client.DeleteItem(context.Background(), &req)
	require.NoError(t, err)
	require.True(t, res.GetSuccess())
}
