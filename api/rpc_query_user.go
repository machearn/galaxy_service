package api

import (
	"context"
	"log"

	"github.com/lib/pq"
	"github.com/machearn/galaxy_service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := server.store.GetMember(ctx, req.GetID())
	if err != nil {
		pqErr := err.(*pq.Error)
		log.Print("cannot get user: ", pqErr)
		return nil, pqErr
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
	}, nil
}
