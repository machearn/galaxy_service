package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/machearn/galaxy_service/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T) Entry {
	member := CreateRandomMember(t)
	item := CreateRandomItem(t)

	arg := CreateEntryParams{
		MemberID:  member.ID,
		ItemID:    item.ID,
		Quantity:  int32(util.GetRandomInt(10, 100)),
		Total:     int32(util.GetRandomInt(10, 1000)),
		CreatedAt: time.Now().UTC(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, member.ID, entry.MemberID)
	require.Equal(t, item.ID, entry.ItemID)
	require.Equal(t, arg.Quantity, entry.Quantity)
	require.Equal(t, arg.Total, entry.Total)
	require.Equal(t, arg.CreatedAt, entry.CreatedAt.UTC())

	require.NotZero(t, entry.ID)

	return entry
}

func CreateRandomEntryByMemberID(t *testing.T, memberID int32) Entry {
	item := CreateRandomItem(t)

	arg := CreateEntryParams{
		MemberID:  memberID,
		ItemID:    item.ID,
		Quantity:  int32(util.GetRandomInt(10, 100)),
		Total:     int32(util.GetRandomInt(10, 1000)),
		CreatedAt: time.Now().UTC(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, memberID, entry.MemberID)
	require.Equal(t, item.ID, entry.ItemID)
	require.Equal(t, arg.Quantity, entry.Quantity)
	require.Equal(t, arg.Total, entry.Total)
	require.Equal(t, arg.CreatedAt, entry.CreatedAt.UTC())

	require.NotZero(t, entry.ID)

	return entry
}

func CreateRandomEntryByItemID(t *testing.T, itemID int32) Entry {
	member := CreateRandomMember(t)

	arg := CreateEntryParams{
		MemberID:  member.ID,
		ItemID:    itemID,
		Quantity:  int32(util.GetRandomInt(10, 100)),
		Total:     int32(util.GetRandomInt(10, 1000)),
		CreatedAt: time.Now().UTC(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, member.ID, entry.MemberID)
	require.Equal(t, itemID, entry.ItemID)
	require.Equal(t, arg.Quantity, entry.Quantity)
	require.Equal(t, arg.Total, entry.Total)
	require.Equal(t, arg.CreatedAt, entry.CreatedAt.UTC())

	require.NotZero(t, entry.ID)

	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := CreateRandomEntry(t)

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry.ID, entry2.ID)
	require.Equal(t, entry.MemberID, entry2.MemberID)
	require.Equal(t, entry.ItemID, entry2.ItemID)
	require.Equal(t, entry.Quantity, entry2.Quantity)
	require.Equal(t, entry.Total, entry2.Total)
	require.Equal(t, entry.CreatedAt.UTC(), entry2.CreatedAt.UTC())

	require.NotZero(t, entry2.ID)
}

func TestListEntries(t *testing.T) {
	var entries []Entry

	for i := 0; i < 10; i++ {
		entries = append(entries, CreateRandomEntry(t))
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 2,
	}

	entries2, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries2)
	require.Len(t, entries2, 5)

	for index, entry := range entries2 {
		require.Equal(t, entries[index].ID, entry.ID)
		require.Equal(t, entries[index].MemberID, entry.MemberID)
		require.Equal(t, entries[index].ItemID, entry.ItemID)
		require.Equal(t, entries[index].Quantity, entry.Quantity)
		require.Equal(t, entries[index].Total, entry.Total)
		require.Equal(t, entries[index].CreatedAt.UTC(), entry.CreatedAt.UTC())

		require.NotZero(t, entry.ID)
	}
}

func TestListEntryByMember(t *testing.T) {
	member := CreateRandomMember(t)

	var entries []Entry

	for i := 0; i < 10; i++ {
		entries = append(entries, CreateRandomEntryByMemberID(t, member.ID))
	}

	arg := ListEntriesByMemberParams{
		MemberID: member.ID,
		Limit:    5,
		Offset:   0,
	}

	entries2, err := testQueries.ListEntriesByMember(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries2)
	require.Len(t, entries2, 5)

	for index, entry := range entries2 {
		require.Equal(t, entries[index].ID, entry.ID)
		require.Equal(t, entries[index].MemberID, entry.MemberID)
		require.Equal(t, entries[index].ItemID, entry.ItemID)
		require.Equal(t, entries[index].Quantity, entry.Quantity)
		require.Equal(t, entries[index].Total, entry.Total)
		require.Equal(t, entries[index].CreatedAt.UTC(), entry.CreatedAt.UTC())

		require.NotZero(t, entry.ID)
	}
}

func TestListEntryByItem(t *testing.T) {
	item := CreateRandomItem(t)

	var entries []Entry

	for i := 0; i < 10; i++ {
		entries = append(entries, CreateRandomEntryByItemID(t, item.ID))
	}

	arg := ListEntriesByItemParams{
		ItemID: item.ID,
		Limit:  5,
		Offset: 0,
	}

	entries2, err := testQueries.ListEntriesByItem(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entries2)
	require.Len(t, entries2, 5)

	for index, entry := range entries2 {
		require.Equal(t, entries[index].ID, entry.ID)
		require.Equal(t, entries[index].MemberID, entry.MemberID)
		require.Equal(t, entries[index].ItemID, entry.ItemID)
		require.Equal(t, entries[index].Quantity, entry.Quantity)
		require.Equal(t, entries[index].Total, entry.Total)
		require.Equal(t, entries[index].CreatedAt.UTC(), entry.CreatedAt.UTC())

		require.NotZero(t, entry.ID)
	}
}

func TestUpdateEntry(t *testing.T) {
	entry := CreateRandomEntry(t)

	arg := UpdateEntryParams{
		ID:        entry.ID,
		MemberID:  sql.NullInt32{Int32: int32(util.GetRandomInt(1, 100)), Valid: true},
		ItemID:    sql.NullInt32{Int32: int32(util.GetRandomInt(1, 100)), Valid: true},
		Quantity:  sql.NullInt32{Int32: int32(util.GetRandomInt(10, 100)), Valid: true},
		Total:     sql.NullInt32{Int32: int32(util.GetRandomInt(10, 1000)), Valid: true},
		CreatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	}

	entry2, err := testQueries.UpdateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, arg.ID, entry2.ID)
	require.Equal(t, arg.MemberID.Int32, entry2.MemberID)
	require.Equal(t, arg.ItemID.Int32, entry2.ItemID)
	require.Equal(t, arg.Quantity.Int32, entry2.Quantity)
	require.Equal(t, arg.Total.Int32, entry2.Total)
	require.Equal(t, arg.CreatedAt.Time.UTC(), entry2.CreatedAt.UTC())

	require.NotZero(t, entry2.ID)
}

func TestPartialUpdateEntry(t *testing.T) {
	entry := CreateRandomEntry(t)

	arg := UpdateEntryParams{
		ID:       entry.ID,
		ItemID:   sql.NullInt32{Int32: int32(util.GetRandomInt(1, 100)), Valid: true},
		Quantity: sql.NullInt32{Int32: int32(util.GetRandomInt(10, 100)), Valid: true},
		Total:    sql.NullInt32{Int32: int32(util.GetRandomInt(10, 1000)), Valid: true},
	}

	entry2, err := testQueries.UpdateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, arg.ID, entry2.ID)
	require.Equal(t, arg.ItemID.Int32, entry2.ItemID)
	require.Equal(t, arg.Quantity.Int32, entry2.Quantity)
	require.Equal(t, arg.Total.Int32, entry2.Total)

	require.NotZero(t, entry2.ID)
}

func TestDeleteEntry(t *testing.T) {
	entry := CreateRandomEntry(t)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)

	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(context.Background(), entry.ID)

	require.Error(t, err)
	require.Empty(t, entry2)
}
