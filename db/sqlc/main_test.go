package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/kvnyijia/bank-app/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries // `type Queries` is defined in db/sqlc/db.go
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Create a new connection to the db for unit tests
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB) // `func New(db DBTX) *Queries` is defined in db/sqlc/db.go

	os.Exit(m.Run())
}
