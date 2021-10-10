package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// discord webhook is a fallback because the
// email notification is our "simulation" pavlok-shock
func webhook(status string) string {

	// TODO instead of env var, retrieve webhook from session db table

	wh := os.Getenv("DISCORD_WEBHOOK")
	av := os.Getenv("GITHUB_AVATAR")
	if wh == "DISABLE_WEBHOOK" {
		return ""
	}

	m := map[string]string {
		"username": "Webhook-twitch-ext-demo",
		"avatar_url": av,
		"content": status,
	}

	buf, err := json.Marshal(m)
	if err != nil {
		log.Print("Marshal failed on webhook")
		log.Print(err)
		return ""
	}

	resp, err := http.Post(wh, "application/json", bytes.NewReader(buf))
	if err != nil {
		log.Print("Post failed for webhook")
		log.Print(err)
		return ""
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Response webhook failed")
		log.Print(err)
		return ""
	}
	log.Printf("Webhook status - %s", data)

	return string(data)
}

