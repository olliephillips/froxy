package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// websocket route handler
func handleSocket(w http.ResponseWriter, r *http.Request) {
	// get alias
	vars := mux.Vars(r)
	accessKey := vars["accesskey"]

	// check sockets allowed
	for _, v := range f.Geofences {
		if v.AccessKey == accessKey {
			if !v.Websocket {
				log.Println("websockets for", accessKey, "are not enabled. Connection denied")
				w.WriteHeader(500)
			}
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// register connection
	froxyConn[accessKey] = true

	// clean up
	defer conn.Close()

	for {
		select {
		case data := <-froxyWS[accessKey]:
			log.Println("websocket data for", accessKey)
			if data != "" {
				err = conn.WriteMessage(websocket.TextMessage, []byte(data))
				if err != nil {
					// unregister conn so we stop sending
					froxyConn[accessKey] = false
					break
				}
			}
		}
	}
}
