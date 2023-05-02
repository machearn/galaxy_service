package api

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/lib/pq"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"github.com/machearn/galaxy_service/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.store.GetMemberByName(ctx, req.GetUsername())
	if err != nil {
		pqErr := err.(*pq.Error)
		log.Print("cannot get user: ", pqErr)
		return nil, pqErr
	}

	if req.GetPassword() != user.Password {
		apiErr := util.NewAPIError("403", "forbiden")
		log.Print("password not match")
		return nil, apiErr
	}

	encryptedAccessToken, accessToken, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("cannot generate access token: ", apiErr)
		return nil, apiErr
	}

	encryptedRefreshToken, refreshToken, err := server.tokenMaker.CreateToken(user.ID, server.config.RefreshTokenDuration)
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("cannot generate refresh token: ", apiErr)
		return nil, apiErr
	}

	uuidString, err := refreshToken.GetJti()
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to get uuid: ", apiErr)
		return nil, apiErr
	}
	ID, err := uuid.Parse(uuidString)
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to parse uuid from string: ", apiErr)
		return nil, apiErr
	}

	issuedAt, err := refreshToken.GetIssuedAt()
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to get issued time: ", apiErr)
		return nil, apiErr
	}
	expiredAt, err := refreshToken.GetExpiration()
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to get expiration time: ", apiErr)
		return nil, apiErr
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
		pqErr := err.(*pq.Error)
		log.Print("cannot create session: ", pqErr)
		return nil, pqErr
	}

	accessExpiredAt, err := accessToken.GetExpiration()
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to get expiration time: ", apiErr)
		return nil, apiErr
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
