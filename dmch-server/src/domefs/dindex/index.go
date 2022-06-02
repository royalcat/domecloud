package dindex

import "go.mongodb.org/mongo-driver/mongo"

type DomeIndex struct {
	VideoInfo *EntryIndex
}

func NewDomeIndex(db *mongo.Database) *DomeIndex {
	return &DomeIndex{
		VideoInfo: NewEntryIndex(db),
	}
}
