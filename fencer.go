package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type fencer struct {
	Data struct {
		Inside bool `json:"inside"`
	} `json:"data"`
	Error interface{} `json:"error"`
}

func callFencer(accessKey string, lat string, lng string) (bool, error) {
	// setup and make call
	endPoint := "https://api.fencer.io/v1.0/position/inside/" + accessKey
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", endPoint, nil)
	req.Header.Add("Authorization", f.Apikey)
	req.Header.Add("Lat-Pos", lat)
	req.Header.Add("Lng-Pos", lng)
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// process response
	data := new(fencer)
	if err := json.Unmarshal(body, &data); err != nil {
		return false, err
	}

	if data.Data.Inside == true {
		return true, nil
	}
	return false, nil
}
