package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	db "github.com/miromax42/wss-chat/db/sqlc"
	app "github.com/miromax42/wss-chat/server"
	"github.com/miromax42/wss-chat/util"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	source := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", config.DBDriver, config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBDatabase)
	conn, err := sql.Open(config.DBDriver, source)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.New(conn)
	server, err := app.New(ctx, config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return server.Srv.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()

		return server.Srv.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
