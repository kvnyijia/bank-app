package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kvnyijia/bank-app/api"
	db "github.com/kvnyijia/bank-app/db/sqlc"

	// _ "github.com/kvnyijia/bank-app/doc/statik"
	mygrpc "github.com/kvnyijia/bank-app/grpc"
	"github.com/kvnyijia/bank-app/pb"
	"github.com/kvnyijia/bank-app/util"
	_ "github.com/lib/pq"

	// "github.com/rakyll/statik/fs"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Info().Msg(">>> Running main.go .......")
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg(">>> cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg(">>> cannot connect to db")
	}

	runDBMigration(config.DBMigration, config.DBSource)

	store := db.NewStore(conn)

	// runGinServer(config, store)

	go runGRPCgatewayServer(config, store)
	runGRPCServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	// HTTP server using Gin
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	log.Info().Msg(">>> Running server .......")
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server:")
	}
}

func runGRPCServer(config util.Config, store db.Store) {
	// gRPC server
	server, err := mygrpc.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg(">>> cannot create server")
	}

	myGrpcLogger := grpc.UnaryInterceptor(mygrpc.GrpcLogger)
	grpcServer := grpc.NewServer(myGrpcLogger)
	pb.RegisterBankAppServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg(">>> cannot create lisener")
	}

	log.Info().Msgf(">>> start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg(">>> cannot start gRPC server")
	}
}

func runGRPCgatewayServer(config util.Config, store db.Store) {
	// gRPC server
	server, err := mygrpc.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg(">>> cannot create server")
	}

	grpcMux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterBankAppHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg(">>> cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// Swagger UI for API docs
	// statikFS, err := fs.New()
	// if err != nil {
	// 	log.Fatal(">>> cannot create statik fs: ", err)
	// }
	// swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	// mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg(">>> cannot create lisener")
	}

	log.Info().Msgf(">>> start gRPC Gateway (HTTP) server at %s", listener.Addr().String())
	mux_with_loggerMiddleware := mygrpc.HttpHandler(mux)
	err = http.Serve(listener, mux_with_loggerMiddleware)
	// err = http.Serve(listener, grpcMux)
	if err != nil {
		log.Fatal().Err(err).Msg(">>> cannot start gRPC Gateway (HTTP) server")
	}
}

func runDBMigration(dbMigration string, dbSource string) {
	migration, err := migrate.New(dbMigration, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg(">>> cannot create new db migration instance")
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg(">>> failed to run migrate up")
	}
	log.Info().Msg(">>> db migrated successfully")
}
