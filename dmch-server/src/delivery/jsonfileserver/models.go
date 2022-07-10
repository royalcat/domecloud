package jsonfileserver

import "dmch-server/src/domefs/dmime"

type Entry struct {
	Name       string         `json:"name"`
	IsListable bool           `json:"isListable"`
	MimeType   dmime.MimeType `json:"mimeType"`
}
