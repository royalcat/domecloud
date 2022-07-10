package domeindex

import (
	"dmch-server/src/domefs/dmime"
	"io"
	"os"
)

type DomeNode interface {
	Name() string
	Stat() (os.FileInfo, error)
	MimeType() dmime.MimeType
}

type DomeFileNode interface {
	DomeNode
	io.Reader
	io.Seeker
}

type DomeListableNode interface {
	DomeNode
	List() ([]DomeNode, error)
}
