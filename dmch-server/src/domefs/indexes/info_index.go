package store

import (
	"github.com/256dpi/lungo"
)

type VideoInfoIndex struct {
	coll lungo.ICollection
}

func NewVideoInfoIndex(db lungo.IDatabase) *VideoInfoIndex {
	infoindex := &VideoInfoIndex{
		coll: db.Collection("info"),
	}
	return infoindex
}
