package delivery

import (
	"dmch-server/src/config"
	"dmch-server/src/delivery/jsonfileserver"
	"dmch-server/src/delivery/player"
	"dmch-server/src/domefs"
	"dmch-server/src/store"
	"encoding/json"
	"net/http"

	"github.com/256dpi/lungo"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type DomeServer struct {
	router http.Handler

	usersStore *store.UsersStore

	domefs *domefs.DomeFS
}

func NewDomeServer(db lungo.IDatabase, usersStore *store.UsersStore) *DomeServer {
	server := &DomeServer{
		router:     httprouter.New(),
		usersStore: usersStore,
		domefs:     domefs.NewDomeFS(db, config.Config.RootFolder, config.Config.CacheFolder),
	}
	server.initRouter()
	return server
}

func (d *DomeServer) initRouter() {
	router := httprouter.New()

	fileserver := jsonfileserver.NewFileServer(d.domefs)
	api := NewApiHandler(d.domefs)
	player := player.NewPlayer(fileserver)

	stripHandler(router,
		"GET", "/file",
		d.AuthWrapper(
			fileserver,
		),
	)

	stripHandler(router,
		"POST", "/file",
		d.AuthWrapper(
			fileserver,
		),
	)

	stripHandler(router,
		"GET", "/api",
		d.AuthWrapper(
			api,
		),
	)

	stripHandler(router,
		"GET", "/player/regtoken",
		d.AuthWrapper(
			http.HandlerFunc(player.GetToken),
		),
	)

	stripHandler(router,
		"GET", "/player/play",
		http.HandlerFunc(player.Play),
	)

	router.Handler(http.MethodGet, "/login", d.AuthWrapper(http.HandlerFunc(d.Login)))

	d.router = cors.AllowAll().Handler(router)
}

func stripHandler(router *httprouter.Router, method, path string, h http.Handler) {
	router.Handler(method, path+"/*path", http.StripPrefix(path, h))
}

func (d *DomeServer) Login(w http.ResponseWriter, r *http.Request) {
	userI := r.Context().Value("user")
	user := userI.(store.User)

	resp, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(resp)
}

func (d *DomeServer) Run() {
	listen := ":5050"
	logrus.Infof("Listening on %s", listen)
	logrus.Error(http.ListenAndServe(listen, d.router))
}
