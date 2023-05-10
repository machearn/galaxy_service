package api

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/machearn/galaxy_service/api_error"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	encryptedAccessToken, accessToken, err := server.tokenMaker.CreateToken(req.GetUserId(), server.config.AccessTokenDuration)
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("cannot generate access token: ", err)
		return nil, apiErr
	}

	encryptedRefreshToken, refreshToken, err := server.tokenMaker.CreateToken(req.GetUserId(), server.config.RefreshTokenDuration)
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
		MemberID:     req.GetUserId(),
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

	return &pb.CreateSessionResponse{
		AccessToken: encryptedAccessToken,
		ExpiredAt:   timestamppb.New(accessExpiredAt),
		Session: &pb.Session{
			ID:           session.ID.String(),
			UserId:       session.MemberID,
			ClientIp:     session.ClientIp,
			UserAgent:    session.UserAgent,
			RefreshToken: session.RefreshToken,
			CreatedAt:    timestamppb.New(session.CreatedAt),
			ExpiredAt:    timestamppb.New(session.ExpiredAt),
		},
	}, nil
}
