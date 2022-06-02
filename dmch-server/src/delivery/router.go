package delivery

import (
	"dmch-server/src/config"
	"dmch-server/src/delivery/jsonfileserver"
	"dmch-server/src/domefs"
	"dmch-server/src/store"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type DomeServer struct {
	router http.Handler

	usersStore *store.UsersStore

	domefs *domefs.DomeFS
}

func NewDomeServer(db *mongo.Database, usersStore *store.UsersStore) *DomeServer {
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

	router.Handler(
		"GET", "/file/*path",
		d.AuthWrapper(
			http.StripPrefix(
				"/file/",
				jsonfileserver.FileServer(d.domefs),
			),
		),
	)

	router.Handler(
		"GET", "/api/*path",
		d.AuthWrapper(
			http.StripPrefix(
				"/api/",
				NewApiHandler(d.domefs),
			),
		),
	)

	router.Handler(http.MethodGet, "/login", d.AuthWrapper(http.HandlerFunc(d.Login)))

	d.router = cors.AllowAll().Handler(router)
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
