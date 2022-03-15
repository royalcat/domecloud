package routes

import (
	"dmch-server/cfs"
	"dmch-server/config"
	"dmch-server/server"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type DmRouter struct {
	router *httprouter.Router
	server *server.DmServer
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
		"GET", "/preview/",
		http.StripPrefix("/preview/*path", http.FileServer(http.Dir(config.Config.Media.PreviewFolder))),
	)
	router.Handler(
		"GET", "/file/",
		http.StripPrefix("/file/*path", http.FileServer(http.FS(cfs.Cfs))),
	)

	const apiPrefix = "/v1/api"
	router.GET(apiPrefix+"/list/*path", d.ListVideo)

	d.router = router
}
