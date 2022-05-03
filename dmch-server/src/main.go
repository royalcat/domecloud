package main

import (
	"context"
	"dmch-server/src/config"
	"dmch-server/src/delivery"
	"dmch-server/src/store"
	"time"

	"github.com/256dpi/lungo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config.Load()
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		logrus.Panicf("Cant connect to MongoDB: %w", err)
	}
	db := client.Database("dome")

	lclient := lungo.MongoClient{client}
	ldb := lclient.Database("dome")

	userStore := store.NewUsersStore(ldb)
	server := delivery.NewDomeServer(db, userStore)
	server.Run()
}
