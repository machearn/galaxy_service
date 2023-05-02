package api

import (
	"context"
	"log"
	"time"

	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"github.com/machearn/galaxy_service/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	arg := db.CreateMemberParams{
		Username:  req.Username,
		Fullname:  req.Fullname,
		Email:     req.Email,
		Password:  req.Password,
		Plan:      req.Plan,
		AutoRenew: req.AutoRenew,
	}

	arg.CreatedAt = time.Now().UTC().Truncate(time.Second)
	arg.ExpiredAt = arg.CreatedAt.Add(util.Plan[arg.Plan]).Truncate(time.Second)

	user, err := server.store.CreateMember(ctx, arg)
	if err != nil {
		log.Print("cannot create user: ", err)
		return nil, err
	}

	return &pb.CreateUserResponse{
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
	}, nil
}
