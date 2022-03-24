package delivery

import (
	"dmch-server/src/cfs"
	"dmch-server/src/delivery/jsonfileserver"
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
		http.StripPrefix("/file/", jsonfileserver.FileServer(cfs.NewDmFS())),
	)

	d.router = router
}

func (d *DmRouter) Run() {
	listen := ":5050"
	logrus.Infof("Listening on %s", listen)
	logrus.Error(http.ListenAndServe(listen, d.router))
}
