package api

import (
	"context"
	"log"

	"github.com/machearn/galaxy_service/pb"
	"github.com/machearn/galaxy_service/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Authorize(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	encryptedAccessToken := req.GetToken()
	token, err := server.tokenMaker.VerifyToken(encryptedAccessToken)
	if err != nil {
		apiErr := util.NewAPIError("403", "forbiden")
		log.Print("cannot verify access token: ", apiErr)
		return nil, apiErr
	}

	ID, err := token.GetJti()
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to get uuid: ", apiErr)
		return nil, apiErr
	}
	var user_id int32
	err = token.Get("member_id", &user_id)
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to get member_id: ", apiErr)
		return nil, apiErr
	}
	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to get issued time: ", apiErr)
		return nil, apiErr
	}
	expiredAt, err := token.GetExpiration()
	if err != nil {
		apiErr := util.NewAPIError("500", "internal error")
		log.Print("failed to get expiration time: ", apiErr)
		return nil, apiErr
	}

	return &pb.AuthResponse{
		ID:        ID,
		UserId:    user_id,
		CreatedAt: timestamppb.New(issuedAt),
		ExpiredAt: timestamppb.New(expiredAt),
	}, nil
}
