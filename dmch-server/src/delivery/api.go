package delivery

import (
	"dmch-server/src/delivery/resthelper"
	"dmch-server/src/domefs"
	"encoding/json"
	"net/http"
	"strings"
)

type ApiHandler struct {
	root *domefs.DomeFS
}

func NewApiHandler(root *domefs.DomeFS) *ApiHandler {
	return &ApiHandler{
		root: root,
	}
}

func (ah *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveMap := map[string]http.Handler{
		"listMedia": http.StripPrefix(
			"listMedia",
			http.HandlerFunc(ah.handleListMedia),
		),
	}

	upaths := strings.Split(r.URL.Path, "/")
	fnname := upaths[0]

	handler, ok := serveMap[fnname]
	if !ok {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	handler.ServeHTTP(w, r)
}

func (d *ApiHandler) handleListMedia(w http.ResponseWriter, r *http.Request) {
	upaths := strings.Split(r.URL.Path, "/")
	upath := "/" + strings.Join(upaths[1:], "/")
	mediaInfos, err := d.root.Api.ListMedia(r.Context(), upath)
	if err != nil {
		resthelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(mediaInfos)
	if err != nil {
		resthelper.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)
}
