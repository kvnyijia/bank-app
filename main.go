package main

import (
	"database/sql"
	"log"

	"github.com/kvnyijia/bank-app/api"
	db "github.com/kvnyijia/bank-app/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver     = "postgres"
	dbSource     = "postgresql://root:110604@localhost:5432/simple_bank?sslmode=disable"
	sererAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(sererAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
