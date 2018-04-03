package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func processWebhooks() {
	// worker to pick up webhooks processing requests from froxyWH channel
	for {
		select {
		case data := <-froxyWH:
			for k, v := range data {
				log.Println("webhook data for", k)
				if v != "" {
					// send hook data on goroutine
					for _, g := range f.Geofences {
						if g.AccessKey == k {
							// get the webhooks
							for _, w := range g.Webhooks {
								go processWebhook(w, v)
							}
						}
					}
				}
			}
		}
	}
}

func processWebhook(hook [3]string, data string) {
	uri := hook[0]
	payload := hook[1]
	event := hook[2]
	status := &Status{}

	if err := json.Unmarshal([]byte(data), &status); err != nil {
		log.Println("failed to unmarshal data")
	}

	// replace hook tokens
	//{client_id}, {geofence_alias}, {lat_pos}, {lng_pos}, {inside}
	var inside string
	if status.Inside == true {
		inside = "true"
	} else {
		inside = "false"
	}

	// first check that this hook should be fired for current status
	if event == "" || event == inside {
		// should be fired process payload
		payload = strings.Replace(payload, "{client_id}", status.ClientID, -1)
		payload = strings.Replace(payload, "{geofence_alias}", status.Alias, -1)
		payload = strings.Replace(payload, "{lat_pos}", status.LatPos, -1)
		payload = strings.Replace(payload, "{lng_pos}", status.LngPos, -1)
		payload = strings.Replace(payload, "{inside}", inside, -1)

		// send request
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		req, err := http.NewRequest("POST", uri, bytes.NewBuffer([]byte(payload)))
		resp, err := client.Do(req)
		if err != nil {
			log.Println("webhook problem", uri, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Println("webhook problem bad response", uri, resp.StatusCode)
		}
	}
}
