package player

import (
	"dmch-server/src/store"
	"time"

	"github.com/bradenaw/juniper/xsync"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerStore struct {
	m xsync.Map[primitive.ObjectID, Access]

	log *logrus.Entry
}

func NewPlayerStore() *PlayerStore {
	return &PlayerStore{
		m:   xsync.Map[primitive.ObjectID, Access]{},
		log: logrus.WithField("service", "player_store"),
	}
}

func (p *PlayerStore) AddAccess(user store.User, path string) Access {
	t := time.Now()
	a := Access{
		Token:        primitive.NewObjectIDFromTimestamp(t),
		Path:         path,
		User:         user,
		Ttl:          time.Hour,
		LastAccessed: t,
	}
	p.m.Store(a.Token, a)
	return a
}

func (p *PlayerStore) GetAccess(token primitive.ObjectID) (Access, bool) {
	a, ok := p.m.Load(token)
	if ok {
		a.LastAccessed = time.Now()
		p.m.Store(a.Token, a)
	}
	return a, ok
}

func (p *PlayerStore) RunCleaner() {
	go func() {
		for now := range time.NewTicker(time.Minute).C {
			p.m.Range(func(key primitive.ObjectID, value Access) bool {
				if value.IsExpired(now) {
					p.m.Delete(key)
				}
				return true
			})
		}
	}()

}

type Access struct {
	Token primitive.ObjectID
	Path  string
	User  store.User

	Ttl          time.Duration
	LastAccessed time.Time
}

func (a *Access) IsExpired(now time.Time) bool {
	return a.LastAccessed.Add(a.Ttl).Before(now)
}
