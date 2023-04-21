package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/machearn/galaxy_service/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomItem(t *testing.T) Item {
	arg := CreateItemParams{
		Name:   util.GetRandomString(5),
		Amount: int32(util.GetRandomInt(10, 1000)),
		Price:  int32(util.GetRandomInt(1000, 10000000)),
	}

	item, err := testQueries.CreateItem(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, item)

	require.Equal(t, arg.Name, item.Name)
	require.Equal(t, arg.Amount, item.Amount)
	require.Equal(t, arg.Price, item.Price)

	require.NotZero(t, item.ID)

	item_row++

	return item
}

func TestCreateItem(t *testing.T) {
	CreateRandomItem(t)
}

func TestGetItem(t *testing.T) {
	item := CreateRandomItem(t)

	item2, err := testQueries.GetItem(context.Background(), item.ID)

	require.NoError(t, err)
	require.NotEmpty(t, item2)

	require.Equal(t, item.Name, item2.Name)
	require.Equal(t, item.Amount, item2.Amount)
	require.Equal(t, item.Price, item2.Price)

	require.NotZero(t, item2.ID)
}

func TestListItems(t *testing.T) {
	var offset int32 = item_row
	var items []Item
	for i := 0; i < 10; i++ {
		item := CreateRandomItem(t)

		items = append(items, item)
	}

	arg := ListItemsParams{
		Limit:  5,
		Offset: offset,
	}

	items2, err := testQueries.ListItems(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, items2)

	require.Len(t, items2, 5)

	for index, item := range items2 {
		require.NotEmpty(t, item)

		require.Equal(t, items[index].Name, item.Name)
		require.Equal(t, items[index].Amount, item.Amount)
		require.Equal(t, items[index].Price, item.Price)

		require.NotZero(t, item.ID)
	}
}

func TestUpdateItem(t *testing.T) {
	item := CreateRandomItem(t)

	arg := UpdateItemParams{
		ID:     item.ID,
		Name:   sql.NullString{String: util.GetRandomString(5), Valid: true},
		Amount: sql.NullInt32{Int32: int32(util.GetRandomInt(10, 1000)), Valid: true},
		Price:  sql.NullInt32{Int32: int32(util.GetRandomInt(1000, 10000000)), Valid: true},
	}

	item2, err := testQueries.UpdateItem(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, item2)

	require.Equal(t, arg.Name.String, item2.Name)
	require.Equal(t, arg.Amount.Int32, item2.Amount)
	require.Equal(t, arg.Price.Int32, item2.Price)

	require.NotZero(t, item2.ID)
}

func TestPatialUpdateItem(t *testing.T) {
	item := CreateRandomItem(t)

	arg := UpdateItemParams{
		ID:     item.ID,
		Amount: sql.NullInt32{Int32: int32(util.GetRandomInt(10, 1000)), Valid: true},
	}

	item2, err := testQueries.UpdateItem(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, item2)

	require.Equal(t, arg.Amount.Int32, item2.Amount)

	require.NotZero(t, item2.ID)
}

func TestDeleteItem(t *testing.T) {
	item := CreateRandomItem(t)

	err := testQueries.DeleteItem(context.Background(), item.ID)

	require.NoError(t, err)

	item2, err := testQueries.GetItem(context.Background(), item.ID)

	require.Error(t, err)
	require.Empty(t, item2)
}
