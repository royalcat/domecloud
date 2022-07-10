package main

import (
	"dmch-server/src/delivery/jsonfileserver"
	"dmch-server/src/domefs/entrymodel"
	"dmch-server/src/store"
	"dmch-server/src/utils/gtoyml"
	"encoding/json"
	"log"
	"os"
	"time"
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

	mi := []entrymodel.MediaInfo{
		{
			MediaType: "video/mp4",
			VideoInfo: &entrymodel.VideoInfo{
				Duration: time.Minute,
				Resolution: entrymodel.Resolution{
					Width:  1920,
					Height: 1080,
				},
			},
		},

		{
			MediaType: "image/png",
			ImageInfo: &entrymodel.ImageInfo{
				Resolution: entrymodel.Resolution{
					Width:  1920,
					Height: 1080,
				},
			},
		},

		{
			MediaType: "audio/flac",
			AudioInfo: &entrymodel.AudioInfo{
				Duration: time.Minute,
			},
		},
	}

	data, _ := json.MarshalIndent(mi, "", "  ")

	f2, err := os.Create("docs/example_media.yml")

	if err != nil {
		log.Fatal(err)
	}

	defer f2.Close()

	_, err = f2.WriteString(string(data))

	if err != nil {
		log.Fatal(err)
	}

}
