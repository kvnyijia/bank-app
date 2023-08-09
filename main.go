package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	db "github.com/kvnyijia/bank-app/db/sqlc"
	mygrpc "github.com/kvnyijia/bank-app/grpc"
	"github.com/kvnyijia/bank-app/pb"
	"github.com/kvnyijia/bank-app/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println(">>> Running main.go .......")
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	// HTTP server using Gin
	// server, err := api.NewServer(config, store)
	// if err != nil {
	// 	log.Fatal("cannot create server:", err)
	// }

	// fmt.Println(">>> Running server .......")
	// err = server.Start(config.ServerAddress)
	// if err != nil {
	// 	log.Fatal("cannot start server:", err)
	// }

	// gRPC server
	server, err := mygrpc.NewServer(config, store)
	if err != nil {
		log.Fatal(">>> cannot create server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBankAppServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal(">>> cannot create lisener")
	}

	log.Printf(">>> start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(">>> cannot start gRPC server")
	}
}
