package main

import (
	"dmch-server/src/delivery/jsonfileserver"
	"dmch-server/src/domefs/entrymodel"
	"dmch-server/src/store"
	"dmch-server/src/utils/gtoyml"
	"log"
	"os"
)

func main() {
	docgen := gtoyml.NewDocGenerator()
	err := docgen.AddModel(entrymodel.MediaInfo{}, jsonfileserver.Entry{}, store.User{})
	if err != nil {
		log.Fatal(err)
	}

	docs, err := docgen.EncodeYaml()

	f, err := os.Create("docs/models.yml")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.WriteString(docs)

	if err != nil {
		log.Fatal(err)
	}

}
