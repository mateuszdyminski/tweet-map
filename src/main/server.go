package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

func LaunchServer(hostname string, port int, staticDir string, r *mux.Router) {
	// Setup static routes
	staticRoutes := make(StaticRoutes, 0)
	if staticDir == "" {
		staticDir = "./"
	}
	staticRoutes = appendStaticRoute(staticRoutes, staticDir)

	// Handle routes
	r.Handle("/{path:.*}", http.FileServer(staticRoutes)).Name("static")

	// Listen on hostname:port
	glog.Infof("Listening on %s:%d...\n", hostname, port)
	http.Handle("/", r)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", hostname, port), nil)
	if err != nil {
		glog.Errorf("Error: %s", err)
	}
}

func appendStaticRoute(sr StaticRoutes, dir string) StaticRoutes {
	if _, err := os.Stat(dir); err != nil {
		glog.Errorf("Error %+v", err)
	}
	return append(sr, http.Dir(dir))
}

type StaticRoutes []http.FileSystem

func (sr StaticRoutes) Open(name string) (f http.File, err error) {
	for _, s := range sr {
		if f, err = s.Open(name); err == nil {
			f = disabledDirListing{f}
			return
		}
	}
	return
}

type disabledDirListing struct {
	http.File
}

func (f disabledDirListing) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}
