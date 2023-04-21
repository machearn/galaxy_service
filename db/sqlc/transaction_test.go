package db

import (
	"context"
	"testing"
	"time"

	"github.com/machearn/galaxy_service/util"
	"github.com/stretchr/testify/require"
)

func TestListEntriesByMember(t *testing.T) {
	store := NewStore(testDB)

	member := CreateRandomMember(t)
	item := CreateRandomItem(t)

	arg := CreateEntryParams{
		MemberID:  member.ID,
		ItemID:    item.ID,
		Quantity:  int32(util.GetRandomInt(1, 1000)),
		Total:     int32(util.GetRandomInt(1, 1000)),
		CreatedAt: time.Now().UTC(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.MemberID, entry.MemberID)
	require.Equal(t, arg.ItemID, entry.ItemID)
	require.Equal(t, arg.Quantity, entry.Quantity)
	require.Equal(t, arg.Total, entry.Total)
	require.Equal(t, arg.CreatedAt, entry.CreatedAt.UTC())

	require.NotZero(t, entry.ID)

	arg2 := ListEntriesByMemberTxParams{
		Username: member.Username,
		Limit:    5,
		Offset:   0,
	}

	entries2, err := store.ListEntriesByMemberTx(context.Background(), arg2)

	require.NoError(t, err)
	require.NotEmpty(t, entries2)
	require.NotEmpty(t, entries2[0])
	require.Len(t, entries2, 1)

	require.Equal(t, entries2[0].ID, entry.ID)
	require.Equal(t, entries2[0].MemberID, entry.MemberID)
	require.Equal(t, entries2[0].ItemID, entry.ItemID)
	require.Equal(t, entries2[0].Quantity, entry.Quantity)
	require.Equal(t, entries2[0].Total, entry.Total)
	require.Equal(t, entries2[0].CreatedAt.UTC(), entry.CreatedAt.UTC())
}
