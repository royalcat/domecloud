package player

import (
	"dmch-server/src/delivery/jsonfileserver"
	"dmch-server/src/store"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Player struct {
	store      *PlayerStore
	fileserver *jsonfileserver.FileServer
}

func NewPlayer(fileserver *jsonfileserver.FileServer) *Player {
	p := &Player{
		store:      NewPlayerStore(),
		fileserver: fileserver,
	}
	p.store.RunCleaner()
	return p
}

func (p *Player) GetToken(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(store.User)
	a := p.store.AddAccess(user, r.URL.Path)

	w.Write([]byte(a.Token.Hex()))
}

func (p *Player) Play(w http.ResponseWriter, r *http.Request) {
	tokenHex := r.URL.Query().Get("token")

	token, err := primitive.ObjectIDFromHex(tokenHex)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	a, ok := p.store.GetAccess(token)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if a.Path != r.URL.Path {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p.fileserver.ServeHTTP(w, r)
	return
}
