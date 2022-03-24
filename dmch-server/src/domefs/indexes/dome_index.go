package store

import "github.com/256dpi/lungo"

type DomeIndex struct {
	VideoInfoIndex *VideoInfoIndex

	engine *lungo.Engine
}

func NewDomeIndex() *DomeIndex {
	opts := lungo.Options{
		Store: lungo.NewMemoryStore(),
	}

	// open database
	client, engine, err := lungo.Open(nil, opts)
	if err != nil {
		panic(err)
	}
	db := client.Database("dome_index")

	return &DomeIndex{
		VideoInfoIndex: NewVideoInfoIndex(db),

		engine: engine,
	}
}

func (d *DomeIndex) Close() {
	if d.engine != nil {
		d.engine.Close()
	}
}
