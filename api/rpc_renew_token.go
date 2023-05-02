package api

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/machearn/galaxy_service/pb"
	"github.com/machearn/galaxy_service/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {
	encryptedRefreshToken := req.GetRefreshToken()

	refreshToken, err := server.tokenMaker.VerifyToken(encryptedRefreshToken)
	if err != nil {
		apiErr := util.NewAPIError("403", "forbiden")
		log.Print("failed to verify refresh token: ", apiErr)
		return nil, apiErr
	}

	uuidString, err := refreshToken.GetJti()
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to get uuid: ", apiErr)
		return nil, apiErr
	}
	id, err := uuid.Parse(uuidString)
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to parse uuid from string: ", apiErr)
		return nil, apiErr
	}

	session, err := server.store.GetSession(ctx, id)
	if err != nil {
		pqErr := err.(*pq.Error)
		log.Print("failed to get session: ", pqErr)
		return nil, pqErr
	}

	if session.IsBlocked {
		apiErr := util.NewAPIError("403", "forbiden")
		log.Print("session is blocked")
		return nil, apiErr
	}

	var memberID int32
	err = refreshToken.Get("member_id", &memberID)
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to get member_id: ", apiErr)
		return nil, apiErr
	}
	if session.MemberID != memberID {
		apiErr := util.NewAPIError("403", "forbiden")
		log.Print("member_id not match")
		return nil, apiErr
	}

	if session.RefreshToken != req.GetRefreshToken() {
		apiErr := util.NewAPIError("403", "forbiden")
		log.Print("refresh token not match")
		return nil, apiErr
	}

	if time.Now().UTC().Truncate(time.Second).After(session.ExpiredAt) {
		apiErr := util.NewAPIError("403", "forbiden")
		log.Print("refresh token expired")
		return nil, apiErr
	}

	encryptedAccessToken, accessToken, err := server.tokenMaker.CreateToken(session.MemberID, server.config.AccessTokenDuration)
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("cannot generate access token: ", apiErr)
		return nil, apiErr
	}
	expiredAt, err := accessToken.GetExpiration()
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to get expiration time: ", apiErr)
		return nil, apiErr
	}

	return &pb.RenewAccessTokenResponse{
		AccessToken: encryptedAccessToken,
		ExpiredAt:   timestamppb.New(expiredAt),
	}, nil
}
