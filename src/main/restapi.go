package main

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

func SetupRestApi(r *mux.Router) {
	api := &RestApi{}
	r.HandleFunc("/restapi/tweets", api.tweets).Methods("POST")
}

type RestApi struct{}

func (r *RestApi) tweets(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("query")

	glog.Infof("Starting search for tweets contain text: %s", query)
	go func() {
		for tweet := range Search(query) {
			tweetToSend := tweet
			glog.Infof("Sending tweet: %s", string(tweetToSend))
			h.broadcast <- tweetToSend
		}
	}()

	json, _ := json.Marshal("")
	w.Write(json)
}
