package api

import (
	"context"
	"log"
	"net/http"

	"github.com/machearn/galaxy_service/api_error"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Authorize(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	encryptedAccessToken := req.GetToken()
	token, err := server.tokenMaker.VerifyToken(encryptedAccessToken)
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusUnauthorized, "unauthorized")
		log.Print("cannot verify access token: ", err)
		return nil, apiErr
	}

	ID, err := token.GetJti()
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to get uuid: ", err)
		return nil, apiErr
	}
	var user_id int32
	err = token.Get("member_id", &user_id)
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to get member_id: ", err)
		return nil, apiErr
	}
	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to get issued time: ", err)
		return nil, apiErr
	}
	expiredAt, err := token.GetExpiration()
	if err != nil {
		apiErr := api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		log.Print("failed to get expiration time: ", err)
		return nil, apiErr
	}

	return &pb.AuthResponse{
		ID:        ID,
		UserId:    user_id,
		CreatedAt: timestamppb.New(issuedAt),
		ExpiredAt: timestamppb.New(expiredAt),
	}, nil
}
