package main

import (
	"dmch-server/src/config"
	"dmch-server/src/delivery"
)

func main() {
	config.Load()

	server := delivery.NewDomeServer()
	server.Run()
}
