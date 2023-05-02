package api

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.store.GetMemberByName(ctx, req.GetUsername())
	if err != nil {
		log.Print("cannot get user: ", err)
		return nil, err
	}

	if req.GetPassword() != user.Password {
		log.Print("password not match")
		return nil, errors.New("password not match")
	}

	encryptedAccessToken, accessToken, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		log.Print("cannot generate access token: ", err)
		return nil, err
	}

	encryptedRefreshToken, refreshToken, err := server.tokenMaker.CreateToken(user.ID, server.config.RefreshTokenDuration)
	if err != nil {
		log.Print("cannot generate refresh token: ", err)
		return nil, err
	}

	uuidString, err := refreshToken.GetJti()
	if err != nil {
		log.Print("failed to get uuid: ", err)
		return nil, err
	}
	ID, err := uuid.Parse(uuidString)
	if err != nil {
		log.Print("failed to parse uuid from string: ", err)
		return nil, err
	}

	issuedAt, err := refreshToken.GetIssuedAt()
	if err != nil {
		log.Print("failed to get issued time: ", err)
		return nil, err
	}
	expiredAt, err := refreshToken.GetExpiration()
	if err != nil {
		log.Print("failed to get expiration time: ", err)
		return nil, err
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           ID,
		MemberID:     user.ID,
		RefreshToken: encryptedRefreshToken,
		ClientIp:     req.GetClientIp(),
		UserAgent:    req.GetUserAgent(),
		CreatedAt:    issuedAt,
		ExpiredAt:    expiredAt,
	})
	if err != nil {
		log.Print("cannot create session: ", err)
		return nil, err
	}

	accessExpiredAt, err := accessToken.GetExpiration()
	if err != nil {
		log.Print("failed to get expiration time: ", err)
		return nil, err
	}

	return &pb.LoginResponse{
		AccessToken:      encryptedAccessToken,
		AccessExpiredAt:  timestamppb.New(accessExpiredAt),
		RefreshToken:     encryptedRefreshToken,
		RefreshExpiredAt: timestamppb.New(expiredAt),
		SessionId:        session.ID.String(),
		User: &pb.User{
			ID:        user.ID,
			Username:  user.Username,
			Fullname:  user.Fullname,
			Email:     user.Email,
			Plan:      user.Plan,
			CreatedAt: timestamppb.New(user.CreatedAt),
			ExpiredAt: timestamppb.New(user.ExpiredAt),
			AutoRenew: user.AutoRenew,
		},
	}, nil
}
