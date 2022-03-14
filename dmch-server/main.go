package main

import (
	"log"
	"net/http"
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

	h := &webdav.Handler{
		Prefix:     "",
		FileSystem: filesystem(),
		LockSystem: webdav.NewMemLS(),
	}

	http.HandleFunc("/", h.ServeHTTP)
	log.Fatal(http.ListenAndServe(":8080", h))
}
