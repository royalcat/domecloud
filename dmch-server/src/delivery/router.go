package delivery

import (
	"dmch-server/src/cfs"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type DmRouter struct {
	router *httprouter.Router
}

func NewDmServer() *DmRouter {
	server := &DmRouter{
		router: httprouter.New(),
	}
	server.initRouter()
	return server
}

func (d *DmRouter) initRouter() {
	router := httprouter.New()

	router.Handler(
		"GET", "/file/*path",
		http.StripPrefix("/file/", http.FileServer(http.FS(cfs.NewDmFS()))),
	)

	d.router = router
}

func (d *DmRouter) Run() {
	listen := ":5050"
	logrus.Infof("Listening on %s", listen)
	logrus.Error(http.ListenAndServe(listen, d.router))
}