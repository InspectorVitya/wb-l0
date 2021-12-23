package main

import (
	"context"
	"errors"
	"github.com/inspectorvitya/wb-l0/internal/application"
	"github.com/inspectorvitya/wb-l0/internal/config"
	"github.com/inspectorvitya/wb-l0/internal/provider/storage/pgsql"
	"github.com/inspectorvitya/wb-l0/internal/provider/streaming"
	"github.com/inspectorvitya/wb-l0/internal/server/httpserver"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := pgsql.New(cfg.DataBase)
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}
	natsStreaming := streaming.New(cfg.Nats)
	app := application.New(db, natsStreaming)
	err = app.Init(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	server := httpserver.New(cfg.HTTP, app)
	go func() {
		log.Println("http server start...")
		if err := server.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("http server http stopped....")
			} else {
				log.Fatalln(err)
			}
		}
	}()

	<-stop
	ctxClose, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = server.Stop(ctxClose)
	if err != nil {
		log.Fatalln(err)
	}
	err = natsStreaming.Stop()
	if err != nil {
		log.Fatalln(err)
	}
}
