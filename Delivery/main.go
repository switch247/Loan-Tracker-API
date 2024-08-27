package main

import (
	Config "Loan-Tracker-API/Config"
	"Loan-Tracker-API/Delivery/routers"
	"context"
	"log"
)

func main() {
	// Connect to the database
	client := Config.ConnectDB()

	// Defer the closing of the database
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	// Setuprouter()
	r := routers.Setuprouter(client)
	if r != nil {
		r.Run(Config.Port)
	} else {
		log.Fatal("Failed to start server")
	}
}
