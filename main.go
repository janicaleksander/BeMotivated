package main

import (
	"github.com/janicaleksander/BeMotivated/api"
	"github.com/janicaleksander/BeMotivated/storage"
	"log"
)

func main() {
	db, err := storage.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	server := api.BuildServer(":5050", db)
	db.Init()
	server.Run()

}
