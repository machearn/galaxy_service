package api

import (
	"context"
	"log"

	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Authorize(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	encryptedAccessToken := req.GetToken()
	token, err := server.tokenMaker.VerifyToken(encryptedAccessToken)
	if err != nil {
		log.Print("cannot verify access token: ", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err.Error())
	}

	ID, err := token.GetJti()
	if err != nil {
		log.Print("failed to get uuid: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}
	var user_id int32
	err = token.Get("member_id", &user_id)
	if err != nil {
		log.Print("failed to get member_id: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}
	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		log.Print("failed to get issued time: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}
	expiredAt, err := token.GetExpiration()
	if err != nil {
		log.Print("failed to get expiration time: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
	}

	return &pb.AuthResponse{
		ID:        ID,
		UserId:    user_id,
		CreatedAt: timestamppb.New(issuedAt),
		ExpiredAt: timestamppb.New(expiredAt),
	}, nil
}
