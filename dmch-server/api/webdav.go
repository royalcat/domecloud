package api

import (
	"context"
	"dmch-server/cfs"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

type WDFs struct {
	cloudfs cfs.Cfs
}

func (wdfs *WDFs) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return wdfs.cloudfs.Mkdir(name, perm)
}

func (wdfs *WDFs) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	return wdfs.cloudfs.Open(name)
}

func (wdfs *WDFs) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return wdfs.cloudfs.Mkdir(name, perm)
}

func WebDav(cloudfs cfs.Cfs) {
	hnd := &webdav.Handler{
		FileSystem: &WDFs{cloudfs: cloudfs},
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				log.Printf("WEBDAV [%s]: %s \n", r.Method, r.URL)
			}
		},
	}
}
