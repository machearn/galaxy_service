package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/machearn/galaxy_service/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomSession(t *testing.T) Session {
	member := CreateRandomMember(t)

	arg := CreateSessionParams{
		ID:           uuid.New(),
		MemberID:     member.ID,
		RefreshToken: util.GetRandomString(32),
		ClientIp:     util.GetRandomString(32),
		UserAgent:    util.GetRandomString(32),
		CreatedAt:    time.Now().UTC().Truncate(time.Second),
		ExpiredAt:    time.Now().AddDate(0, 0, 30).UTC().Truncate(time.Second),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.MemberID, session.MemberID)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.CreatedAt, session.CreatedAt.UTC())
	require.Equal(t, arg.ExpiredAt, session.ExpiredAt.UTC())
	require.False(t, session.IsBlocked)

	return session
}

func TestCreateSession(t *testing.T) {
	CreateRandomSession(t)
}

func TestGetSession(t *testing.T) {
	session := CreateRandomSession(t)

	session2, err := testQueries.GetSession(context.Background(), session.ID)

	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session.ID, session2.ID)
	require.Equal(t, session.MemberID, session2.MemberID)
	require.Equal(t, session.RefreshToken, session2.RefreshToken)
	require.Equal(t, session.ClientIp, session2.ClientIp)
	require.Equal(t, session.UserAgent, session2.UserAgent)
	require.Equal(t, session.CreatedAt, session2.CreatedAt)
	require.Equal(t, session.ExpiredAt, session2.ExpiredAt)
	require.Equal(t, session.IsBlocked, session2.IsBlocked)
}
