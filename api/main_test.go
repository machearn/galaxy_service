package api

import (
	"net"
	"os"
	"testing"

	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	token "github.com/machearn/galaxy_service/token/maker"
	"github.com/machearn/galaxy_service/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startTestServer(t *testing.T, store db.Store, tokenMaker token.TokenMaker) net.Listener {
	config, err := util.LoadConfig("../")
	require.NoError(t, err)
	require.NotEmpty(t, config)

	if tokenMaker == nil {
		tokenMaker, err = token.NewTokenMaker(config.TokenSymmetricKey)
		require.NoError(t, err)
		require.NotEmpty(t, tokenMaker)
	}

	server, err := NewServer(config, store, tokenMaker, nil)
	require.NoError(t, err)

	grpcServer := grpc.NewServer()
	pb.RegisterGalaxyServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return listener
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
