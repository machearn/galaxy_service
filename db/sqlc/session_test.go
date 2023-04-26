package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/machearn/galaxy_service/util"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	member := CreateRandomMember(t)

	arg := CreateSessionParams{
		ID:        uuid.New(),
		MemberID:  member.ID,
		Token:     util.GetRandomString(32),
		ClientIp:  util.GetRandomString(32),
		UserAgent: util.GetRandomString(32),
		IsActive:  true,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		ExpiredAt: time.Now().AddDate(0, 0, 30).UTC().Truncate(time.Second),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.MemberID, session.MemberID)
	require.Equal(t, arg.Token, session.Token)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.IsActive, session.IsActive)
	require.Equal(t, arg.CreatedAt, session.CreatedAt.UTC())
	require.Equal(t, arg.ExpiredAt, session.ExpiredAt.UTC())
}
