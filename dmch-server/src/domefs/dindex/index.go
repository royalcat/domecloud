package dindex

import "go.mongodb.org/mongo-driver/mongo"

type DomeIndex struct {
	VideoInfo *VideoInfoIndex
}

func NewDomeIndex(db *mongo.Database) *DomeIndex {
	return &DomeIndex{
		VideoInfo: NewVideoInfoIndex(db),
	}
}
