package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	db "github.com/miromax42/wss-chat/db/sqlc"
	"github.com/miromax42/wss-chat/server"
	"github.com/miromax42/wss-chat/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	source := fmt.Sprintf("%s://%s:%s@$%s:%d/%s?sslmode=disable", config.DBDriver, config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBDatabase)
	conn, err := sql.Open(config.DBDriver, source)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.New(conn)
	server, err := server.New(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerPort)
	if err == nil {
		log.Fatal("caanot start server: ", err)
	}
}
