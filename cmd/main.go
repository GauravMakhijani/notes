package main

import (
	"log"

	"github.com/GauravMakhijani/notes/internal/database"
	"github.com/GauravMakhijani/notes/internal/service"
	"github.com/urfave/negroni"
)

func main() {
	store := database.NewStore()
	if err := store.AutoMigrate(); err != nil {
		log.Println(err)
		panic("Failed to migrate database")
	}

	service := service.NewService(store)
	appRouter := initRouter(service)
	server := negroni.Classic()
	server.UseHandler(appRouter)
	log.Println("Starting server on port 8080")
	server.Run(":8080")
	return
}
