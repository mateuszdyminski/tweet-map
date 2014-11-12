package main

import (
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

func SetupWebSocket(r *mux.Router, redisConn string) {
	api := &WebSocketApi{}

	// setup connection to redis
	api.setupRedisConn(redisConn)

	// run websocket hub
	go h.run()

	// run redis pub-sub to make WS connections the same on all nodes
	go api.subscribeToTopic()

	r.HandleFunc("/wsapi/ws", api.serveWs)
}

type WebSocketApi struct {
	redisPool *redis.Pool
}

// serverWs handles websocket requests from the peer.
func (r *WebSocketApi) serveWs(w http.ResponseWriter, req *http.Request) {
	glog.Infof("Registering client to WS")
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		glog.Errorf("Error %+v", err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump()
}

func (r *WebSocketApi) subscribeToTopic() {
	conn := r.redisPool.Get()
	defer conn.Close()

	psc := redis.PubSubConn{conn}
	psc.Subscribe("ws")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			glog.Infof("%s: message: %s\n", v.Channel, v.Data)
		case redis.Subscription:
			glog.Infof("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			glog.Errorf("Error: %v", v)
		}
	}
}

func (r *WebSocketApi) setupRedisConn(redisConn string) {
	glog.Infof("Starting redis connection pool to: %s", redisConn)
	r.redisPool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisConn)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
