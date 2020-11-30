package main

import (
	"flag"
	"github.com/kompir/app-pointment/server"
	"github.com/kompir/app-pointment/server/repositories"
	"github.com/kompir/app-pointment/server/services"
	"log"
	"os"
	"syscall"
)

var (
	addrFlag        = flag.String("addr", ":8008", "HTTP server address")
	notifierURIFlag = flag.String("notifier", "http://localhost:9000", "Notifier API URI")
	dbFlag          = flag.String("db", "db.json", "Path to db.json file")
	dbCfgFlag       = flag.String("db-cfg", ".db.config.json", "Path to .db.config.json file")
)

func main() {
	flag.Parse()
	db := repositories.NewDB(*dbFlag, *dbCfgFlag)
	repo := repositories.NewReminders(db)
	service := services.NewReminders(repo)
	backend := server.New(*addrFlag, service)
	saver := services.NewSaver(service)
	notifier := services.NewNotifier(*notifierURIFlag, service)
	if err := db.Start(); err != nil {
		log.Fatalf("could not start file database service: %v", err)
	}
	go saver.Start()
	go notifier.Start()
	go func() {
		if err := backend.Start(); err != nil {
			log.Fatalf("could not start backend api service: %v", err)
		}
	}()
	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	server.ListenForSignals(signals, backend, saver, notifier, db)
}