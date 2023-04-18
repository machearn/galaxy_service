package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/machearn/galaxy/util"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var row_index int32 = 0

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("failed to load config:", err)
		os.Exit(1)
	}

	dbConn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
		os.Exit(1)
	}

	testQueries = New(dbConn)

	os.Exit(m.Run())
}
