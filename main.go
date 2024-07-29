package main

import (
	"fmt"
	"github.com/janicaleksander/BeMotivated/api"
	"github.com/janicaleksander/BeMotivated/storage"
	"os"
)

func main() {
	db, err := storage.NewPostgresDB()
	if err != nil {
		fmt.Println("xd")
	}

	port := os.Getenv("PORT")

	server := api.BuildServer(":"+port, db)
	db.Init()
	server.Run()

}
