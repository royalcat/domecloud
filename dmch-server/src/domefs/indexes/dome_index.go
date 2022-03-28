package indexes

type DomeIndex struct {
	VideoInfo *VideoInfoIndex
}

func NewDomeIndex() *DomeIndex {
	return &DomeIndex{
		VideoInfo: NewVideoInfoIndex(),
	}
}
