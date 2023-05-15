package api

import (
	"context"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"github.com/machearn/galaxy_service/util"
	"github.com/machearn/galaxy_service/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	callback := func(member db.Member) error {
		taskPayload := worker.EmailPayload{
			Username: member.Username,
		}
		opts := []asynq.Option{
			asynq.MaxRetry(3),
			asynq.Timeout(10 * time.Second),
			asynq.Queue("critical"),
		}
		return server.distributor.DistributeTaskSendVerificationEmail(ctx, taskPayload, opts...)
	}

	user, err := server.store.CreateMemberTx(ctx, arg, callback)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Print("cannot create user, database error: ", pqErr)
			if pqErr.Code[:2] == "23" {
				return nil, status.Errorf(codes.InvalidArgument, "invalid input: %v", err.Error())
			}
			return nil, status.Errorf(codes.Internal, "internal database error: %v", err.Error())
		}
		log.Print("cannot create user: ", err)
		return nil, status.Errorf(codes.Internal, "internal error: %v", err.Error())
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
