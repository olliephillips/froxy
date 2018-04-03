package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Status struct {
	ClientID string `json:"client_id"`
	Alias    string `json:"geofence_alias"`
	LatPos   string `json:"lat_pos"`
	LngPos   string `json:"lng_pos"`
	Inside   bool   `json:"inside"`
}

func client(w http.ResponseWriter, r *http.Request) {
	// get alias
	vars := mux.Vars(r)
	accessKey := vars["accesskey"]
	lat := r.Header.Get("Lat-Pos")
	lng := r.Header.Get("Lng-Pos")
	clientID := r.Header.Get("Client-ID")

	_, ok := froxyStatus[accessKey][clientID]
	if !ok {
		mutex.Lock()
		froxyStatus[accessKey][clientID] = false
		mutex.Unlock()
	}

	// make request
	inside, err := callFencer(accessKey, lat, lng)
	if err != nil {
		log.Println(err)
	}

	// check map
	if status, ok := froxyStatus[accessKey][clientID]; ok {
		if status != inside {
			// need to package this as JSON with some other info
			webhooks := false
			status := &Status{
				ClientID: clientID,
				Inside:   inside,
				LatPos:   lat,
				LngPos:   lng,
			}

			// get fence alias
			for _, v := range f.Geofences {
				if v.AccessKey == accessKey {
					status.Alias = v.Alias
					// check web if web hooks set up
					if v.Webhooks != nil {
						webhooks = true
					}
				}
			}

			// marshal to JSON and put on channel
			json, _ := json.Marshal(status)

			// changed better appraise any socket connections (if we have any)
			if froxyConn[accessKey] {
				froxyWS[accessKey] <- string(json)
			} else {
				log.Println("no websocket connection on", accessKey)
			}

			// put on channel for webhooks
			if webhooks {
				data := make(map[string]string)
				data[accessKey] = string(json)
				froxyWH <- data
			} else {
				log.Println("no webhooks defined for", accessKey)
			}
		}
	}

	// update map
	mutex.Lock()
	froxyStatus[accessKey][clientID] = inside
	mutex.Unlock()

	// really just for debugging
	var output string
	if inside {
		output = "true"
	} else {
		output = "false"
	}
	fmt.Fprintf(w, output)
}
