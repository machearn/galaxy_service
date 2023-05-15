package api

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {
	encryptedRefreshToken := req.GetRefreshToken()

	refreshToken, err := server.tokenMaker.VerifyToken(encryptedRefreshToken)
	if err != nil {
		log.Print("failed to verify refresh token: ", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err.Error())
	}

	uuidString, err := refreshToken.GetJti()
	if err != nil {
		log.Print("failed to get uuid: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}
	id, err := uuid.Parse(uuidString)
	if err != nil {
		log.Print("failed to parse uuid from string: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	session, err := server.store.GetSession(ctx, id)
	if err != nil {
		log.Print("failed to get session: ", err)
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "session not found: %v", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	if session.IsBlocked {
		log.Print("session is blocked")
		return nil, status.Errorf(codes.Unauthenticated, "session is blocked")
	}

	var memberID int32
	err = refreshToken.Get("member_id", &memberID)
	if err != nil {
		log.Print("failed to get member_id: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}
	if session.MemberID != memberID {
		log.Print("member_id not match")
		return nil, status.Errorf(codes.Unauthenticated, "user not match")
	}

	if session.RefreshToken != req.GetRefreshToken() {
		log.Print("refresh token not match")
		return nil, status.Errorf(codes.Unauthenticated, "refresh token not match")
	}

	if time.Now().UTC().Truncate(time.Second).After(session.ExpiredAt) {
		log.Print("refresh token expired")
		return nil, status.Errorf(codes.Unauthenticated, "refresh token expired")
	}

	encryptedAccessToken, accessToken, err := server.tokenMaker.CreateToken(session.MemberID, server.config.AccessTokenDuration)
	if err != nil {
		log.Print("cannot generate access token: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}
	expiredAt, err := accessToken.GetExpiration()
	if err != nil {
		log.Print("failed to get expiration time: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	return &pb.RenewAccessTokenResponse{
		AccessToken: encryptedAccessToken,
		ExpiredAt:   timestamppb.New(expiredAt),
	}, nil
}
