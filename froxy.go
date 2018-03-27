package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
)

type froxy struct {
	Apikey    string     `toml:"apikey"`
	Geofences []geofence `toml:"geofence"`
}

type geofence struct {
	Alias     string   `toml:"alias"`
	AccessKey string   `toml:"accesskey"`
	Websocket bool     `toml:"websocket"`
	Webhooks  []string `toml:"webhooks"`
}

var f froxy
var mutex = &sync.Mutex{}

var froxyWS = make(map[string]chan string)
var froxyStatus = make(map[string]map[string]bool)
var froxyConn = make(map[string]bool)

func main() {
	// load config.toml
	if _, err := toml.DecodeFile("config.toml", &f); err != nil {
		log.Fatal(err)
	}

	// init tracking map and message channel for each geofence
	for _, v := range f.Geofences {
		froxyStatus[v.AccessKey] = make(map[string]bool)
		froxyWS[v.AccessKey] = make(chan string)
	}

	// router
	r := mux.NewRouter()

	// application route handlers
	// client applications
	r.HandleFunc(`/client/{accesskey:[a-zA-Z0-9=\-\/]+}`, client).Methods("GET")

	// websockets
	r.HandleFunc(`/ws/{accesskey:[a-zA-Z0-9=\-\/]+}`, handleSocket).Methods("GET")

	// start server
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:9000",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	// start server
	log.Println("Starting server on port 9000")
	log.Fatal(srv.ListenAndServe())
}
