package api

import (
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	token "github.com/machearn/galaxy_service/token/maker"
	"github.com/machearn/galaxy_service/util"
	"github.com/machearn/galaxy_service/worker"
)

type Server struct {
	pb.UnimplementedGalaxyServer
	config      util.Config
	store       db.Store
	tokenMaker  token.TokenMaker
	distributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, maker token.TokenMaker, distributor worker.TaskDistributor) (*Server, error) {
	return &Server{
		config:      config,
		store:       store,
		tokenMaker:  maker,
		distributor: distributor,
	}, nil
}
