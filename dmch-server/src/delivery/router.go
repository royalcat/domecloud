package delivery

import (
	"dmch-server/src/delivery/jsonfileserver"
	"dmch-server/src/domefs"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type DomeServer struct {
	router *httprouter.Router
}

func NewDomeServer() *DomeServer {
	server := &DomeServer{
		router: httprouter.New(),
	}
	server.initRouter()
	return server
}

func (d *DomeServer) initRouter() {
	router := httprouter.New()

	router.Handler(
		"GET", "/file/*path",
		http.StripPrefix("/file/", jsonfileserver.FileServer(domefs.NewDomeFS())),
	)

	d.router = router
}

func (d *DomeServer) Run() {
	listen := ":5050"
	logrus.Infof("Listening on %s", listen)
	logrus.Error(http.ListenAndServe(listen, d.router))
}
