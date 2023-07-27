package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kvnyijia/bank-app/api"
	db "github.com/kvnyijia/bank-app/db/sqlc"
	"github.com/kvnyijia/bank-app/util"
	_ "github.com/lib/pq"
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
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	fmt.Println(">>> Running server .......")
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
