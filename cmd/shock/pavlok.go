package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

)


// call pavlok API 
func callAPI(token string) string {
	log.Printf("API access - %s", token)

	base := "https://app.pavlok.com/api/v1/stimuli/shock/50"

	m := map[string]string{"access_token": token}

	buf, err := json.Marshal(m)
	if err != nil {
		log.Print("Marshal failed on token")
		log.Print(err)
		return ""
	}
	log.Printf("Sending token - %s", buf)

	resp, err := http.Post(base, "application/json", bytes.NewReader(buf))
	if err != nil {
		log.Printf("Post failed - %s: ", base)
		log.Print(err)
		return ""
	}
	defer resp.Body.Close()

	// TODO limit reader and timeout
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Response failed")
		log.Print(err)
		return ""
	}
	log.Printf("API status - %s", data)

	return string(data)
}


