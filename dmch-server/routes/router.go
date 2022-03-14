package routes

import (
	"github.com/julienschmidt/httprouter"
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

func (d *DmRouter) initRouter() {}
