package main

import (
	"dmch-server/src/config"
	"dmch-server/src/delivery"
	"log"
	"os"

	"golang.org/x/net/webdav"
)

func filesystem() webdav.FileSystem {

	if err := os.Mkdir("./data", os.ModePerm); !os.IsExist(err) {
		log.Fatalf("FATAL %v", err)
	}
	log.Printf("INFO using local filesystem at %s", "./data")
	return webdav.Dir("./data")

}

func main() {

	config.Load()
	router := delivery.NewDmServer()
	router.Run()
}