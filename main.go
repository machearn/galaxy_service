package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/machearn/galaxy_service/api"
	db "github.com/machearn/galaxy_service/db/sqlc"
	"github.com/machearn/galaxy_service/mail"
	"github.com/machearn/galaxy_service/pb"
	token "github.com/machearn/galaxy_service/token/maker"
	"github.com/machearn/galaxy_service/util"
	"github.com/machearn/galaxy_service/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = RunMigration(config.MigrateURL, config.DBSource)
	if err != nil {
		log.Fatal("cannot run migration:", err)
	}

	store := db.NewStore(conn)

	tokenMaker, err := token.NewTokenMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker: ", err)
	}

	redisOpt := asynq.RedisClientOpt{Addr: config.RedisAddress}

	go RunTaskProcessor(config, store, redisOpt)

	distributor := worker.NewRedisTaskDistributor(redisOpt)

	server, err := api.NewServer(config, store, tokenMaker, distributor)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGalaxyServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	log.Printf("server listening at %s", config.GRPCServerAddress)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func RunMigration(url string, dbSource string) error {
	migrate, err := migrate.New(url, dbSource)
	if err != nil {
		log.Printf("failed to create migration: %s", err.Error())
		return err
	}

	if err := migrate.Up(); err != nil {
		log.Printf("failed to run migration: %s", err.Error())
		return err
	}
	return nil
}

func RunTaskProcessor(config util.Config, store db.Store, opt asynq.RedisClientOpt) error {
	sender := mail.NewGmailSender("Machearn", config.TestGmailAddress, config.TestGmailPassword)
	processor := worker.NewRedisTaskProcessor(store, sender, opt)
	err := processor.Start()
	if err != nil {
		log.Printf("failed to start processor: %s", err.Error())
		return err
	}

	return processor.Start()
}
