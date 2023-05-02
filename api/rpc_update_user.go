package api

import (
	"context"
	"database/sql"
	"log"
	"time"

	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"github.com/machearn/galaxy_service/util"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	arg := db.UpdateMemberParams{
		ID: req.GetID(),
	}

	if req.Username != nil {
		arg.Username = sql.NullString{
			String: req.GetUsername(),
			Valid:  true,
		}
	}

	if req.Fullname != nil {
		arg.Fullname = sql.NullString{
			String: req.GetFullname(),
			Valid:  true,
		}
	}

	if req.Email != nil {
		arg.Email = sql.NullString{
			String: req.GetEmail(),
			Valid:  true,
		}
	}

	if req.Password != nil {
		arg.Password = sql.NullString{
			String: req.GetPassword(),
			Valid:  true,
		}
	}

	if req.Plan != nil {
		arg.Plan = sql.NullInt32{
			Int32: req.GetPlan(),
			Valid: true,
		}
		arg.ExpiredAt = sql.NullTime{
			Time:  time.Now().UTC().Add(util.Plan[req.GetPlan()]).Truncate(time.Second),
			Valid: true,
		}
	}

	user, err := server.store.UpdateMember(ctx, arg)
	if err != nil {
		log.Print("cannot update user: ", err)
		return nil, err
	}

	rep := pb.UpdateUserResponse{
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
	}

	return &rep, nil
}
