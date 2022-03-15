package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (s *DmRouter) ListVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := context.Background()

	dpath := strings.TrimPrefix(r.URL.Path, "/list")
	info, err := s.server.ListVideos(ctx, dpath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
