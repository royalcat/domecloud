package main

import (
	"context"
	"dmch-server/src/config"
	"dmch-server/src/delivery"
	"dmch-server/src/store"
	"io/fs"
	"time"

	"github.com/256dpi/lungo"
)

func main() {
	config.Load()
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	// if err != nil {
	// 	logrus.Panicf("Cant connect to MongoDB: %w", err)
	// }
	// db := client.Database("dome")
	//lclient := lungo.MongoClient{client}

	opts := lungo.Options{
		Store: lungo.NewFileStore("lungo.db", fs.ModePerm),
	}
	// open database
	lclient, engine, err := lungo.Open(ctx, opts)
	if err != nil {
		panic(err)
	}
	defer engine.Close()

	ldb := lclient.Database("dome")

	userStore := store.NewUsersStore(ldb)
	server := delivery.NewDomeServer(ldb, userStore)
	server.Run()
}
