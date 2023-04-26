package api

import (
	"net"
	"os"
	"testing"

	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/pb"
	"github.com/machearn/galaxy_service/util"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func startTestServer(t *testing.T, store db.Store) net.Listener {
	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	server, err := NewServer(config, store)
	require.NoError(t, err)

	grpcServer := grpc.NewServer()
	pb.RegisterGalaxyServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.ServerAddress)
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return listener
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
