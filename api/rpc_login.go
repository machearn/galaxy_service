package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/machearn/galaxy_service/api_error"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.store.GetMemberByName(ctx, req.GetUsername())
	if err != nil {
		var apiErr *api_error.APIError
		if err == sql.ErrNoRows {
			apiErr = api_error.NewAPIError(http.StatusBadRequest, "user not found")
		} else {
			apiErr = api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		}
		log.Print("cannot get user: ", err)
		return nil, apiErr
	}

	if req.GetPassword() != user.Password {
		apiErr := api_error.NewAPIError(http.StatusUnauthorized, "password not match")
		log.Print("password not match")
		return nil, apiErr
	}

	encryptedAccessToken, accessToken, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("cannot generate access token: ", err)
		return nil, apiErr
	}

	encryptedRefreshToken, refreshToken, err := server.tokenMaker.CreateToken(user.ID, server.config.RefreshTokenDuration)
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("cannot generate refresh token: ", err)
		return nil, apiErr
	}

	uuidString, err := refreshToken.GetJti()
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to get uuid: ", err)
		return nil, apiErr
	}
	ID, err := uuid.Parse(uuidString)
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to parse uuid from string: ", err)
		return nil, apiErr
	}

	issuedAt, err := refreshToken.GetIssuedAt()
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to get issued time: ", err)
		return nil, apiErr
	}
	expiredAt, err := refreshToken.GetExpiration()
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to get expiration time: ", err)
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
		var apiErr *api_error.APIError
		if pqErr.Code[:2] == "23" {
			apiErr = api_error.NewAPIError(http.StatusBadRequest, "invalid input")
		} else {
			apiErr = api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		}
		log.Print("cannot create session: ", err)
		return nil, apiErr
	}

	accessExpiredAt, err := accessToken.GetExpiration()
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to get expiration time: ", err)
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
