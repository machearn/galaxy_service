package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/machearn/galaxy_service/api_error"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {
	encryptedRefreshToken := req.GetRefreshToken()

	refreshToken, err := server.tokenMaker.VerifyToken(encryptedRefreshToken)
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusUnauthorized, "unauthorized")
		log.Print("failed to verify refresh token: ", err)
		return nil, apiErr
	}

	uuidString, err := refreshToken.GetJti()
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to get uuid: ", err)
		return nil, apiErr
	}
	id, err := uuid.Parse(uuidString)
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to parse uuid from string: ", err)
		return nil, apiErr
	}

	session, err := server.store.GetSession(ctx, id)
	if err != nil {
		var apiErr *api_error.APIError
		if err == sql.ErrNoRows {
			apiErr = api_error.NewAPIError(http.StatusNotFound, "session not found")
		} else {
			apiErr = api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		}
		log.Print("failed to get session: ", err)
		return nil, apiErr
	}

	if session.IsBlocked {
		apiErr := api_error.NewAPIError(http.StatusUnauthorized, "internal error")
		log.Print("session is blocked")
		return nil, apiErr
	}

	var memberID int32
	err = refreshToken.Get("member_id", &memberID)
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to get member_id: ", err)
		return nil, apiErr
	}
	if session.MemberID != memberID {
		apiErr := api_error.NewAPIError(http.StatusUnauthorized, "user not match")
		log.Print("member_id not match")
		return nil, apiErr
	}

	if session.RefreshToken != req.GetRefreshToken() {
		apiErr := api_error.NewAPIError(http.StatusUnauthorized, "refresh token not match")
		log.Print("refresh token not match")
		return nil, apiErr
	}

	if time.Now().UTC().Truncate(time.Second).After(session.ExpiredAt) {
		apiErr := api_error.NewAPIError(http.StatusUnauthorized, "refresh token expired")
		log.Print("refresh token expired")
		return nil, apiErr
	}

	encryptedAccessToken, accessToken, err := server.tokenMaker.CreateToken(session.MemberID, server.config.AccessTokenDuration)
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("cannot generate access token: ", err)
		return nil, apiErr
	}
	expiredAt, err := accessToken.GetExpiration()
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to get expiration time: ", apiErr)
		return nil, apiErr
	}

	return &pb.RenewAccessTokenResponse{
		AccessToken: encryptedAccessToken,
		ExpiredAt:   timestamppb.New(expiredAt),
	}, nil
}
