package delivery

import (
	"context"
	"dmch-server/src/store"
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

		if !isPathAccessable(r.URL.Path, *user) {
			httpWriteError(w, http.StatusForbidden, errors.New("Forbidden, access denied"))
			return
		}

		r = r.WithContext(context.WithValue(ctx, "user", *user))

		h.ServeHTTP(w, r)
	})
}

func httpWriteError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

func isPathAccessable(path string, user store.User) bool {
	parts := strings.Split(strings.TrimLeft(path, "/"), "/")
	if parts == nil || len(parts) == 0 {
		return false
	}

	switch parts[0] {
	case "login":
		return true
	case "user":
		if len(parts) >= 2 {
			return user.Username == parts[1]
		}
	case "share":
		if len(parts) >= 2 {
			switch parts[1] {
			case "all":
				return true
			case "user":
				return user.Username == parts[2]
			}
		}
	case "external":
		return false // FIXME
	}
	return false
}
