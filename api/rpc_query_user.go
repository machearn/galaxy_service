package api

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/machearn/galaxy_service/api_error"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := server.store.GetMember(ctx, req.GetID())
	if err != nil {
		var apiErr *api_error.APIError
		if err == sql.ErrNoRows {
			apiErr = api_error.NewAPIError(http.StatusNotFound, "user not found")
		} else {
			apiErr = api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		}
		log.Print("cannot get user: ", err)
		return nil, apiErr
	}

	return &pb.GetUserResponse{
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
		Password: user.Password,
	}, nil
}

func (server *Server) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserResponse, error) {
	user, err := server.store.GetMemberByName(ctx, req.GetUsername())
	if err != nil {
		var apiErr *api_error.APIError
		if err == sql.ErrNoRows {
			apiErr = api_error.NewAPIError(http.StatusNotFound, "user not found")
		} else {
			apiErr = api_error.NewAPIError(http.StatusInternalServerError, "internal error")
		}
		log.Print("cannot get user: ", err)
		return nil, apiErr
	}

	return &pb.GetUserResponse{
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
		Password: user.Password,
	}, nil
}
