package api

import (
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"github.com/machearn/galaxy_service/util"
)

type Server struct {
	pb.UnimplementedGalaxyServer
	config util.Config
	store  db.Store
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	return &Server{
		config: config,
		store:  store,
	}, nil
}
