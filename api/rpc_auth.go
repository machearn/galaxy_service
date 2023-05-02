package api

import (
	"context"
	"log"

	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Authorize(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	encryptedAccessToken := req.GetToken()
	token, err := server.tokenMaker.VerifyToken(encryptedAccessToken)
	if err != nil {
		log.Print("cannot verify access token: ", err)
		return nil, err
	}

	ID, err := token.GetJti()
	if err != nil {
		log.Print("failed to get uuid: ", err)
		return nil, err
	}
	var user_id int32
	err = token.Get("member_id", &user_id)
	if err != nil {
		log.Print("failed to get member_id: ", err)
		return nil, err
	}
	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		log.Print("failed to get issued time: ", err)
		return nil, err
	}
	expiredAt, err := token.GetExpiration()
	if err != nil {
		log.Print("failed to get expiration time: ", err)
		return nil, err
	}

	return &pb.AuthResponse{
		ID:        ID,
		UserId:    user_id,
		CreatedAt: timestamppb.New(issuedAt),
		ExpiredAt: timestamppb.New(expiredAt),
	}, nil
}
