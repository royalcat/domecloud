package delivery

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

func (d *DomeServer) AuthWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			httpWriteError(w, http.StatusUnauthorized, errors.New("Unauthorized, invalid auth header"))
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		username := pair[0]
		password := pair[1]

		if len(pair) != 2 {
			httpWriteError(w, http.StatusUnauthorized, errors.New("Unauthorized, invalid auth header"))
			return
		}

		user, err := d.usersStore.GetUser(ctx, username)
		if err != nil {
			httpWriteError(w, http.StatusInternalServerError, err)
			return
		}
		if user == nil {
			httpWriteError(w, http.StatusUnauthorized, errors.New("Unauthorized, user not found"))
			return
		}

		if user.Password != password { // TODO: check hash, not plain
			httpWriteError(w, http.StatusUnauthorized, errors.New("Unauthorized, invalid password"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user", *user))

		h.ServeHTTP(w, r)
	})
}

func httpWriteError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
