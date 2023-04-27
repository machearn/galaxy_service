package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
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
		CreatedAt: time.Now().UTC().Truncate(time.Second),
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

func TestCreateMemberTx(t *testing.T) {
	store := NewStore(testDB)

	arg1 := CreateMemberTxParams{
		Username:  util.GetRandomString(5),
		Fullname:  util.GetRandomString(10),
		Email:     util.GetRandomString(5) + "@" + util.GetRandomString(5) + ".com",
		Password:  util.GetRandomString(10),
		Plan:      1,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		ExpiredAt: time.Now().UTC().Truncate(time.Second),
		AutoRenew: true,
	}

	arg2 := CreateSessionTxParams{
		ID:        uuid.New(),
		Token:     util.GetRandomString(10),
		ClientIp:  util.GetRandomString(10),
		UserAgent: util.GetRandomString(10),
		IsActive:  true,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		ExpiredAt: time.Now().UTC().Truncate(time.Second),
	}

	result, err := store.CreateMemberTx(context.Background(), arg1, arg2)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, arg1.Username, result.member.Username)
	require.Equal(t, arg1.Fullname, result.member.Fullname)
	require.Equal(t, arg1.Email, result.member.Email)
	require.Equal(t, arg1.Password, result.member.Password)
	require.Equal(t, arg1.Plan, result.member.Plan)
	require.Equal(t, arg1.CreatedAt, result.member.CreatedAt.UTC())
	require.Equal(t, arg1.ExpiredAt, result.member.ExpiredAt.UTC())
	require.Equal(t, arg1.AutoRenew, result.member.AutoRenew)

	require.NotZero(t, result.member.ID)
	require.Equal(t, arg2.ID, result.session.ID)
	require.Equal(t, result.member.ID, result.session.MemberID)
	require.Equal(t, arg2.Token, result.session.Token)
	require.Equal(t, arg2.ClientIp, result.session.ClientIp)
	require.Equal(t, arg2.UserAgent, result.session.UserAgent)
	require.Equal(t, arg2.IsActive, result.session.IsActive)
	require.Equal(t, arg2.CreatedAt, result.session.CreatedAt.UTC())
	require.Equal(t, arg2.ExpiredAt, result.session.ExpiredAt.UTC())
}
