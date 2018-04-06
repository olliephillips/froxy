package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type froxy struct {
	Apikey    string     `toml:"apikey"`
	Geofences []geofence `toml:"geofence"`
}

type geofence struct {
	Alias     string      `toml:"alias"`
	AccessKey string      `toml:"accesskey"`
	Websocket bool        `toml:"websocket"`
	Webhooks  [][3]string `toml:"webhooks"`
}

var f froxy
var mutex = &sync.Mutex{}

var froxyWS = make(map[string]chan string)
var froxyWH = make(chan map[string]string)
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

	// cors policy
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	// configure server
	srv := &http.Server{
		Handler:      c.Handler(r),
		Addr:         ":9000",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	// start webhooks worker in goroutine
	go processWebhooks()

	// start server
	log.Println("Starting server on port 9000")
	log.Fatal(srv.ListenAndServe())
}
