package api

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/lib/pq"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	encryptedAccessToken, accessToken, err := server.tokenMaker.CreateToken(req.GetUserId(), server.config.AccessTokenDuration)
	if err != nil {
		log.Print("cannot generate access token: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	encryptedRefreshToken, refreshToken, err := server.tokenMaker.CreateToken(req.GetUserId(), server.config.RefreshTokenDuration)
	if err != nil {
		log.Print("cannot generate refresh token: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	uuidString, err := refreshToken.GetJti()
	if err != nil {
		log.Print("failed to get uuid: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}
	ID, err := uuid.Parse(uuidString)
	if err != nil {
		log.Print("failed to parse uuid from string: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	issuedAt, err := refreshToken.GetIssuedAt()
	if err != nil {
		log.Print("failed to get issued time: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}
	expiredAt, err := refreshToken.GetExpiration()
	if err != nil {
		log.Print("failed to get expiration time: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
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
		log.Print("cannot create session: ", err)
		pqErr := err.(*pq.Error)
		if pqErr.Code[:2] == "23" {
			return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	accessExpiredAt, err := accessToken.GetExpiration()
	if err != nil {
		log.Print("failed to get expiration time: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
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
