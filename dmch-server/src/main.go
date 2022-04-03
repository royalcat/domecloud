package main

import (
	"dmch-server/src/config"
	"dmch-server/src/delivery"
	"dmch-server/src/store"

	"github.com/256dpi/lungo"
)

func main() {
	config.Load()

	opts := lungo.Options{
		Store: lungo.NewFileStore("./db.bson", 0666),
	}

	// open database
	client, _, err := lungo.Open(nil, opts)
	if err != nil {
		panic(err)
	}
	db := client.Database("dome")

	userStore := store.NewUsersStore(db)
	server := delivery.NewDomeServer(userStore)
	server.Run()
}
