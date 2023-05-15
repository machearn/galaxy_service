package api

import (
	"context"
	"database/sql"
	"log"

	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := server.store.GetMember(ctx, req.GetID())
	if err != nil {
		log.Print("cannot get user: ", err)
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
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
		log.Print("cannot get user: ", err)
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
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
