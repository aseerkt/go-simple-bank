package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/aseerkt/go-simple-bank/pkg/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config := utils.LoadConfig("../..")

	var err error
	testDB, err = sql.Open(config.DBDriver, config.DBUrl)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	defer testDB.Close()

	testQueries = New(testDB)

	os.Exit(m.Run())
}
