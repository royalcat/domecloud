package delivery

import (
	"dmch-server/src/delivery/jsonfileserver"
	"dmch-server/src/domefs"
	"dmch-server/src/store"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type DomeServer struct {
	router *httprouter.Router

	usersStore *store.UsersStore
}

func NewDomeServer(usersStore *store.UsersStore) *DomeServer {
	server := &DomeServer{
		router:     httprouter.New(),
		usersStore: usersStore,
	}
	server.initRouter()
	return server
}

func (d *DomeServer) initRouter() {
	router := httprouter.New()

	router.Handler(
		"GET", "/file/*path",
		d.AuthWrapper(
			http.StripPrefix(
				"/file/",
				jsonfileserver.FileServer(domefs.NewDomeFS()),
			),
		),
	)

	d.router = router
}

func (d *DomeServer) Run() {
	listen := ":5050"
	logrus.Infof("Listening on %s", listen)
	logrus.Error(http.ListenAndServe(listen, d.router))
}
