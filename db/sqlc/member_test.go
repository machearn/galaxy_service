package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/machearn/galaxy/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomRow(t *testing.T) Member {
	arg := CreateMemberParams{
		Username:  util.GetRandomString(5),
		Fullname:  util.GetRandomString(10),
		Email:     util.GetRandomString(10) + "@gmail.com",
		Plan:      int32(util.GetRandomInt(1, 3)),
		CreatedAt: time.Now().UTC(),
		ExpiredAt: time.Now().AddDate(0, 0, 30).UTC(),
		AutoRenew: true,
	}

	member, err := testQueries.CreateMember(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, member)

	require.Equal(t, arg.Username, member.Username)
	require.Equal(t, arg.Fullname, member.Fullname)
	require.Equal(t, arg.Email, member.Email)
	require.Equal(t, arg.Plan, member.Plan)
	require.Equal(t, arg.CreatedAt, member.CreatedAt.UTC())
	require.Equal(t, arg.ExpiredAt, member.ExpiredAt.UTC())
	require.Equal(t, arg.AutoRenew, member.AutoRenew)

	require.NotZero(t, member.ID)

	row_index++

	return member
}

func TestCreatMember(t *testing.T) {
	CreateRandomRow(t)
}

func TestGetMember(t *testing.T) {
	member := CreateRandomRow(t)

	member2, err := testQueries.GetMember(context.Background(), member.ID)

	require.NoError(t, err)
	require.NotEmpty(t, member2)

	require.Equal(t, member.Username, member2.Username)
	require.Equal(t, member.Fullname, member2.Fullname)
	require.Equal(t, member.Email, member2.Email)
	require.Equal(t, member.Plan, member2.Plan)
	require.Equal(t, member.CreatedAt.UTC(), member2.CreatedAt.UTC())
	require.Equal(t, member.ExpiredAt.UTC(), member2.ExpiredAt.UTC())
	require.Equal(t, member.AutoRenew, member2.AutoRenew)

	require.NotZero(t, member2.ID)
}

func TestListMembers(t *testing.T) {
	var members []Member
	offset := row_index
	for i := 0; i < 10; i++ {
		member := CreateRandomRow(t)

		members = append(members, member)
	}

	args := ListMembersParams{
		Limit:  5,
		Offset: offset,
	}

	members2, err := testQueries.ListMembers(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, members2, 5)

	for index, member := range members2 {
		require.NotEmpty(t, member)

		require.Equal(t, members[index].Username, member.Username)
		require.Equal(t, members[index].Fullname, member.Fullname)
		require.Equal(t, members[index].Email, member.Email)
		require.Equal(t, members[index].Plan, member.Plan)
		require.Equal(t, members[index].CreatedAt.UTC(), member.CreatedAt.UTC())
		require.Equal(t, members[index].ExpiredAt.UTC(), member.ExpiredAt.UTC())
		require.Equal(t, members[index].AutoRenew, member.AutoRenew)

		require.NotZero(t, member.ID)
	}
}

func TestUpdateMember(t *testing.T) {
	member := CreateRandomRow(t)

	arg2 := UpdateMemberParams{
		ID:        member.ID,
		Username:  sql.NullString{util.GetRandomString(5), true},
		Fullname:  sql.NullString{util.GetRandomString(10), true},
		Email:     sql.NullString{util.GetRandomString(10) + "@gmail.com", true},
		Plan:      sql.NullInt32{int32(util.GetRandomInt(1, 3)), true},
		ExpiredAt: sql.NullTime{time.Now().AddDate(0, 0, 30).UTC(), true},
		AutoRenew: sql.NullBool{false, true},
	}

	member2, err := testQueries.UpdateMember(context.Background(), arg2)

	require.NoError(t, err)
	require.NotEmpty(t, member2)

	require.Equal(t, arg2.Username.String, member2.Username)
	require.Equal(t, arg2.Fullname.String, member2.Fullname)
	require.Equal(t, arg2.Email.String, member2.Email)
	require.Equal(t, arg2.Plan.Int32, member2.Plan)
	require.Equal(t, arg2.ExpiredAt.Time.UTC(), member2.ExpiredAt.UTC())
	require.Equal(t, arg2.AutoRenew.Bool, member2.AutoRenew)
}

func TestPartionUpdateMember(t *testing.T) {
	member := CreateRandomRow(t)

	arg2 := UpdateMemberParams{
		ID:       member.ID,
		Username: sql.NullString{util.GetRandomString(5), true},
		Email:    sql.NullString{util.GetRandomString(10) + "@gmail.com", true},
		Plan:     sql.NullInt32{int32(util.GetRandomInt(1, 3)), true},
	}

	member2, err := testQueries.UpdateMember(context.Background(), arg2)

	require.NoError(t, err)
	require.NotEmpty(t, member2)

	require.Equal(t, arg2.Username.String, member2.Username)
	require.Equal(t, arg2.Email.String, member2.Email)
	require.Equal(t, arg2.Plan.Int32, member2.Plan)
}

func TestDeleteMember(t *testing.T) {
	member := CreateRandomRow(t)

	err := testQueries.DeleteMember(context.Background(), member.ID)

	require.NoError(t, err)

	member2, err := testQueries.GetMember(context.Background(), member.ID)

	require.Error(t, err)
	require.Empty(t, member2)
}
