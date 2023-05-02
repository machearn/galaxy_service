package api

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {
	encryptedRefreshToken := req.GetRefreshToken()

	refreshToken, err := server.tokenMaker.VerifyToken(encryptedRefreshToken)
	if err != nil {
		log.Print("failed to verify refresh token: ", err)
		return nil, err
	}

	uuidString, err := refreshToken.GetJti()
	if err != nil {
		log.Print("failed to get uuid: ", err)
		return nil, err
	}
	id, err := uuid.Parse(uuidString)
	if err != nil {
		log.Print("failed to parse uuid from string: ", err)
		return nil, err
	}

	session, err := server.store.GetSession(ctx, id)
	if err != nil {
		log.Print("failed to get session: ", err)
		return nil, err
	}

	if session.IsBlocked {
		log.Print("session is blocked")
		return nil, err
	}

	var memberID int32
	err = refreshToken.Get("member_id", &memberID)
	if err != nil {
		log.Print("failed to get member_id: ", err)
		return nil, err
	}
	if session.MemberID != memberID {
		log.Print("member_id not match")
		return nil, err
	}

	if session.RefreshToken != req.GetRefreshToken() {
		log.Print("refresh token not match")
		return nil, errors.New("refresh token not match")
	}

	if time.Now().UTC().Truncate(time.Second).After(session.ExpiredAt) {
		log.Print("refresh token expired")
		return nil, errors.New("refresh token expired")
	}

	encryptedAccessToken, accessToken, err := server.tokenMaker.CreateToken(session.MemberID, server.config.AccessTokenDuration)
	if err != nil {
		log.Print("cannot generate access token: ", err)
		return nil, err
	}
	expiredAt, err := accessToken.GetExpiration()
	if err != nil {
		log.Print("failed to get expiration time: ", err)
		return nil, err
	}

	return &pb.RenewAccessTokenResponse{
		AccessToken: encryptedAccessToken,
		ExpiredAt:   timestamppb.New(expiredAt),
	}, nil
}
