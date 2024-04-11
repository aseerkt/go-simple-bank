package main

import (
	"database/sql"
	"log"

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

	store := db.NewStore(conn)

	server := api.NewServer(store, &config)

	server.LoadRoutes()

	server.Start(config.ServerAddress)
}
