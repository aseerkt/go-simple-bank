package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/aseerkt/go-simple-bank/pkg/api"
	"github.com/aseerkt/go-simple-bank/pkg/db"
	"github.com/aseerkt/go-simple-bank/pkg/utils"

	_ "github.com/lib/pq"
)

func main() {
	config := utils.LoadConfig(".")

	conn, err := sql.Open(config.DBDriver, config.DBUrl)

	if err != nil {
		log.Fatal("unable to connect to db: ", err)
	}

	runDBMigrations(config.MigrateUrl, config.DBUrl)

	store := db.NewStore(conn)

	server := api.NewServer(store, &config)

	server.LoadRoutes()

	server.Start(config.ServerAddress)
}

func runDBMigrations(migratePath string, dbUrl string) {
	migrateInstance, err := migrate.New(migratePath, dbUrl)

	if err != nil {
		log.Fatal("cannot create new migrate instance", err)
	}

	if err := migrateInstance.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up", err)
	}

	log.Println("db migrate successfully")
}
