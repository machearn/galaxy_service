package api

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/machearn/galaxy_service/db/mock"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	mock_token "github.com/machearn/galaxy_service/token/mock"
	"github.com/machearn/galaxy_service/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type eqCreateMemberParamsMatcher struct {
	arg db.CreateMemberParams
}

type eqCreateSessionParamsMatcher struct {
	arg db.CreateSessionParams
}

type eqDurationMatcher struct {
	arg util.Config
}

func (e eqCreateMemberParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateMemberParams)
	if !ok {
		return false
	}

	if e.arg.Username != arg.Username {
		return false
	}
	if e.arg.Fullname != arg.Fullname {
		return false
	}
	if e.arg.Email != arg.Email {
		return false
	}
	if e.arg.Password != arg.Password {
		return false
	}
	if e.arg.Plan != arg.Plan {
		return false
	}
	if e.arg.AutoRenew != arg.AutoRenew {
		return false
	}

	return true
}

func (e eqCreateMemberParamsMatcher) String() string {
	return "matches with eqCreateMemberParamsMatcher"
}

func (e eqCreateSessionParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateSessionParams)
	if !ok {
		return false
	}

	if e.arg.MemberID != arg.MemberID {
		return false
	}
	if e.arg.ClientIp != arg.ClientIp {
		return false
	}
	if e.arg.UserAgent != arg.UserAgent {
		return false
	}

	return true
}

func (e eqCreateSessionParamsMatcher) String() string {
	return "matches with eqCreateSessionParamsMatcher"
}

func (e eqDurationMatcher) Matches(x interface{}) bool {
	arg, ok := x.(time.Duration)
	if !ok {
		return false
	}

	if e.arg.AccessTokenDuration == arg {
		return true
	}
	if e.arg.RefreshTokenDuration == arg {
		return true
	}

	return false
}

func (e eqDurationMatcher) String() string {
	return "matches with eqDurationMatcher"
}

func EqCreateMemberParams(arg db.CreateMemberParams) gomock.Matcher {
	return eqCreateMemberParamsMatcher{arg: arg}
}

func EqCreateSessionParams(arg db.CreateSessionParams) gomock.Matcher {
	return eqCreateSessionParamsMatcher{arg: arg}
}

func EqDuration(arg util.Config) gomock.Matcher {
	return eqDurationMatcher{arg: arg}
}

func creatRandomUser() db.Member {
	plan := int32(util.GetRandomInt(1, 3))
	user := db.Member{
		ID:        int32(util.GetRandomInt(1, 100)),
		Username:  util.GetRandomString(10),
		Fullname:  util.GetRandomString(10),
		Email:     util.GetRandomString(10) + "@gmail.com",
		Password:  util.GetRandomString(10),
		Plan:      plan,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		ExpiredAt: time.Now().UTC().Add(util.Plan[plan]).Truncate(time.Second),
		AutoRenew: true,
	}

	return user
}

func TestCreateUserAPI(t *testing.T) {
	user := creatRandomUser()

	req := pb.CreateUserRequest{
		Username:  user.Username,
		Fullname:  user.Fullname,
		Email:     user.Email,
		Password:  user.Password,
		Plan:      user.Plan,
		AutoRenew: user.AutoRenew,
	}

	arg := db.CreateMemberParams{
		Username:  req.Username,
		Fullname:  req.Fullname,
		Email:     req.Email,
		Password:  req.Password,
		Plan:      req.Plan,
		CreatedAt: user.CreatedAt,
		ExpiredAt: user.ExpiredAt,
		AutoRenew: req.AutoRenew,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().CreateMemberTx(gomock.Any(), EqCreateMemberParams(arg), gomock.Any()).Times(1).Return(user, nil)

	listener := startTestServer(t, store, nil)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	res, err := client.CreateUser(context.Background(), &req)
	require.NoError(t, err)
	require.Equal(t, user.ID, res.GetUser().ID)
	require.Equal(t, user.Username, res.GetUser().Username)
	require.Equal(t, user.Fullname, res.GetUser().Fullname)
	require.Equal(t, user.Email, res.GetUser().Email)
	require.Equal(t, user.Plan, res.GetUser().Plan)
	require.Equal(t, user.AutoRenew, res.GetUser().AutoRenew)
}

func TestGetUserByUsernameAPI(t *testing.T) {
	user := creatRandomUser()

	req := pb.GetUserByUsernameRequest{
		Username: user.Username,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().GetMemberByName(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)

	listener := startTestServer(t, store, nil)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	res, err := client.GetUserByUsername(context.Background(), &req)
	require.NoError(t, err)
	require.Equal(t, user.ID, res.GetUser().ID)
	require.Equal(t, user.Username, res.GetUser().Username)
	require.Equal(t, user.Fullname, res.GetUser().Fullname)
	require.Equal(t, user.Email, res.GetUser().Email)
	require.Equal(t, user.Plan, res.GetUser().Plan)
	require.Equal(t, user.CreatedAt, res.GetUser().CreatedAt.AsTime())
	require.Equal(t, user.ExpiredAt, res.GetUser().ExpiredAt.AsTime())
	require.Equal(t, user.AutoRenew, res.GetUser().AutoRenew)
}

func TestGetUserAPI(t *testing.T) {
	user := creatRandomUser()

	req := pb.GetUserRequest{
		ID: user.ID,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().GetMember(gomock.Any(), gomock.Eq(user.ID)).Times(1).Return(user, nil)

	listener := startTestServer(t, store, nil)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	res, err := client.GetUser(context.Background(), &req)
	require.NoError(t, err)
	require.Equal(t, user.ID, res.GetUser().ID)
	require.Equal(t, user.Username, res.GetUser().Username)
	require.Equal(t, user.Fullname, res.GetUser().Fullname)
	require.Equal(t, user.Email, res.GetUser().Email)
	require.Equal(t, user.Plan, res.GetUser().Plan)
	require.Equal(t, user.CreatedAt, res.GetUser().CreatedAt.AsTime())
	require.Equal(t, user.ExpiredAt, res.GetUser().ExpiredAt.AsTime())
	require.Equal(t, user.AutoRenew, res.GetUser().AutoRenew)
}

func TestLoginAPI(t *testing.T) {
	config, err := util.LoadConfig("../")
	require.NoError(t, err)
	require.NotEmpty(t, config)

	user := creatRandomUser()

	req := pb.LoginRequest{
		Username:  user.Username,
		Password:  user.Password,
		ClientIp:  "1.1.1.1",
		UserAgent: util.GetRandomString(32),
	}

	arg := db.CreateSessionParams{
		MemberID:  user.ID,
		ClientIp:  req.ClientIp,
		UserAgent: req.UserAgent,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	maker := mock_token.NewMockTokenMaker(ctrl)

	token := paseto.NewToken()
	now := time.Now().UTC().Truncate(time.Second)
	exp := now.Add(config.AccessTokenDuration).Truncate(time.Second)
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(exp)
	ID := uuid.New()
	token.SetJti(ID.String())
	token.Set("member_id", user.ID)
	key := paseto.NewV4SymmetricKey()
	encryptedToken := token.V4Encrypt(key, nil)
	maker.EXPECT().CreateToken(gomock.Eq(user.ID), EqDuration(config)).Times(2).Return(encryptedToken, &token, nil)

	session := db.Session{
		ID:           ID,
		MemberID:     user.ID,
		RefreshToken: encryptedToken,
		ClientIp:     req.ClientIp,
		UserAgent:    req.UserAgent,
		IsBlocked:    false,
		CreatedAt:    now,
		ExpiredAt:    exp,
	}

	store.EXPECT().GetMemberByName(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
	store.EXPECT().CreateSession(gomock.Any(), EqCreateSessionParams(arg)).Times(1).Return(session, nil)

	listener := startTestServer(t, store, maker)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	res, err := client.Login(context.Background(), &req)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, encryptedToken, res.GetAccessToken())
	require.Equal(t, exp, res.GetAccessExpiredAt().AsTime())
	require.Equal(t, encryptedToken, res.GetRefreshToken())
	require.Equal(t, exp, res.GetRefreshExpiredAt().AsTime())
	require.Equal(t, ID.String(), res.GetSessionId())
	require.Equal(t, user.ID, res.GetUser().ID)
	require.Equal(t, user.Username, res.GetUser().Username)
	require.Equal(t, user.Fullname, res.GetUser().Fullname)
	require.Equal(t, user.Email, res.GetUser().Email)
	require.Equal(t, user.Plan, res.GetUser().Plan)
	require.Equal(t, user.CreatedAt, res.GetUser().CreatedAt.AsTime())
	require.Equal(t, user.ExpiredAt, res.GetUser().ExpiredAt.AsTime())
	require.Equal(t, user.AutoRenew, res.GetUser().AutoRenew)
}

func TestUpdateUserAPI(t *testing.T) {
	user := creatRandomUser()

	newUsername := util.GetRandomString(12)

	req := pb.UpdateUserRequest{
		ID:       user.ID,
		Username: &newUsername,
	}

	arg := db.UpdateMemberParams{
		ID:       user.ID,
		Username: sql.NullString{String: newUsername, Valid: true},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user.Username = newUsername

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().UpdateMember(gomock.Any(), gomock.Eq(arg)).Times(1).Return(user, nil)

	listener := startTestServer(t, store, nil)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	res, err := client.UpdateUser(context.Background(), &req)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, user.ID, res.GetUser().ID)
	require.Equal(t, user.Username, res.GetUser().Username)
}

func TestValidAuthAPI(t *testing.T) {
	user := creatRandomUser()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	token := paseto.NewToken()
	now := time.Now().UTC().Truncate(time.Second)
	exp := now.Add(time.Hour).Truncate(time.Second)
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(exp)
	ID := uuid.New()
	token.SetJti(ID.String())
	token.Set("member_id", user.ID)
	key := paseto.NewV4SymmetricKey()
	encryptedToken := token.V4Encrypt(key, nil)

	maker := mock_token.NewMockTokenMaker(ctrl)
	maker.EXPECT().VerifyToken(gomock.Eq(encryptedToken)).Times(1).Return(&token, nil)

	listener := startTestServer(t, nil, maker)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	req := pb.AuthRequest{
		Token: encryptedToken,
	}
	res, err := client.Authorize(context.Background(), &req)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, ID.String(), res.GetID())
	require.Equal(t, user.ID, res.GetUserId())
	require.Equal(t, now, res.GetCreatedAt().AsTime())
	require.Equal(t, exp, res.GetExpiredAt().AsTime())
}

func TestInvalidAuthAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	encryptedToken := util.GetRandomString(32)
	maker := mock_token.NewMockTokenMaker(ctrl)
	maker.EXPECT().VerifyToken(gomock.Eq(encryptedToken)).Times(1).Return(nil, errors.New("invalid token"))

	listener := startTestServer(t, nil, maker)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	req := pb.AuthRequest{
		Token: encryptedToken,
	}
	res, err := client.Authorize(context.Background(), &req)
	require.Error(t, err)
	require.Empty(t, res)
}

func TestRenewTokenAPI(t *testing.T) {
	config, err := util.LoadConfig("../")
	require.NoError(t, err)
	require.NotEmpty(t, config)

	user := creatRandomUser()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	maker := mock_token.NewMockTokenMaker(ctrl)

	token := paseto.NewToken()
	now := time.Now().UTC().Truncate(time.Second)
	exp := now.Add(config.AccessTokenDuration).Truncate(time.Second)
	token.SetIssuedAt(now)
	token.SetNotBefore(now)
	token.SetExpiration(exp)
	ID := uuid.New()
	token.SetJti(ID.String())
	token.Set("member_id", user.ID)
	key := paseto.NewV4SymmetricKey()
	encryptedToken := token.V4Encrypt(key, nil)
	maker.EXPECT().VerifyToken(gomock.Eq(encryptedToken)).Times(1).Return(&token, nil)
	maker.EXPECT().CreateToken(gomock.Eq(user.ID), gomock.Eq(config.AccessTokenDuration)).Times(1).Return(encryptedToken, &token, nil)

	session := db.Session{
		ID:           ID,
		MemberID:     user.ID,
		RefreshToken: encryptedToken,
		ClientIp:     "1.1.1.1",
		UserAgent:    util.GetRandomString(32),
		IsBlocked:    false,
		CreatedAt:    now,
		ExpiredAt:    exp,
	}

	store.EXPECT().GetSession(gomock.Any(), gomock.Eq(ID)).Times(1).Return(session, nil)

	listener := startTestServer(t, store, maker)
	defer listener.Close()

	conn, err := grpc.Dial(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewGalaxyClient(conn)

	req := pb.RenewAccessTokenRequest{
		RefreshToken: encryptedToken,
	}
	res, err := client.RenewAccessToken(context.Background(), &req)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, encryptedToken, res.GetAccessToken())
	require.Equal(t, exp, res.GetExpiredAt().AsTime())
}
