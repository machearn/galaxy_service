package token

import (
	"testing"
	"time"

	"github.com/machearn/galaxy_service/util"
	"github.com/stretchr/testify/require"
)

func TestNewTokenMaker(t *testing.T) {
	key := util.GetRandomString(32)

	maker, err := NewTokenMaker(key)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	memberID := int32(util.GetRandomInt(1, 20))
	duration := time.Minute

	issuedAt := time.Now().UTC().Truncate(time.Second)
	expiredAt := issuedAt.Add(duration).Truncate(time.Second)

	signed, token, err := maker.CreateToken(memberID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, signed)
	require.NotEmpty(t, token)

	parsedToken, err := maker.VerifyToken(signed)

	require.NoError(t, err)
	require.NotEmpty(t, parsedToken)

	var parsedMemberID int32
	parsedToken.Get("member_id", &parsedMemberID)
	require.Equal(t, memberID, parsedMemberID)
	parsedIssue, err := parsedToken.GetIssuedAt()
	require.NoError(t, err)
	require.NotEmpty(t, parsedIssue)
	parsedExpire, err := parsedToken.GetExpiration()
	require.NoError(t, err)
	require.NotEmpty(t, parsedExpire)
	require.WithinDuration(t, issuedAt, parsedIssue, time.Second)
	require.WithinDuration(t, expiredAt, parsedExpire, time.Second)
}

func TestExpiredToken(t *testing.T) {
	key := util.GetRandomString(32)

	maker, err := NewTokenMaker(key)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	memberID := int32(util.GetRandomInt(1, 20))
	duration := time.Second

	signed, token, err := maker.CreateToken(memberID, duration)
	require.NoError(t, err)
	require.NotEmpty(t, signed)
	require.NotEmpty(t, token)

	time.Sleep(time.Second * 2)

	parsedToken, err := maker.VerifyToken(signed)

	require.Error(t, err)
	require.Empty(t, parsedToken)
}
