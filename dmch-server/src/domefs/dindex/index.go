package dindex

import "github.com/256dpi/lungo"

type DomeIndex struct {
	VideoInfo *EntryIndex
}

func NewDomeIndex(db lungo.IDatabase) *DomeIndex {
	return &DomeIndex{
		VideoInfo: NewEntryIndex(db),
	}
}
