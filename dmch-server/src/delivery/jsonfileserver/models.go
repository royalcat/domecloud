package jsonfileserver

type Entry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}
