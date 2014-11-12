package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

var (
	staticDir string
	hostname  string
	port      int
	redisConn string
)

func init() {
	// Flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s -dir [static_dir]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&hostname, "h", "localhost", "hostname")
	flag.StringVar(&staticDir, "dir", "app", "app directory")
	flag.StringVar(&redisConn, "redisConn", "localhost:6379", "hostname:port of redis")
	flag.IntVar(&port, "p", 8080, "port")
}

func main() {
	// start web server with as many cpu as possible
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse flags
	flag.Parse()

	// HTTP routing
	r := mux.NewRouter()

	// start rest api
	SetupRestApi(r)

	// start websocket
	SetupWebSocket(r, redisConn)

	// start file server
	LaunchServer(hostname, port, staticDir, r)

	glog.Infof("Done")
}
