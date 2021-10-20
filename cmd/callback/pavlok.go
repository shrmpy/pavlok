package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int    `json:"created_at"`
	User         struct {
		ID         int         `json:"id"`
		Gender     interface{} `json:"gender"`
		Age        interface{} `json:"age"`
		Location   interface{} `json:"location"`
		Developer  bool        `json:"developer"`
		Onboarded  bool        `json:"onboarded"`
		Locked     bool        `json:"locked"`
		PictureURL interface{} `json:"picture_url"`
	} `json:"user"`
	Device     interface{} `json:"device"`
	UniqueHash string      `json:"unique_hash"`
}

// obtain jwt by trading auth code
func tradeToken(ev events.APIGatewayProxyRequest) string {
	code := ev.QueryStringParameters["code"]
	log.Printf("Code to be exchanged - %s", code)

	cid := os.Getenv("PAVLOK_CLIENT_ID")
	sec := os.Getenv("PAVLOK_CLIENT_SECRET")
	redir := os.Getenv("PAVLOK_REDIRECT_URI")
	base := "https://app.pavlok.com/oauth/token"

	m := map[string]string{
		"client_id":     cid,
		"client_secret": sec,
		"redirect_uri":  redir,
		"grant_type":    "authorization_code",
		"code":          code,
	}

	buf, err := json.Marshal(m)
	if err != nil {
		log.Printf("Marshal failed for map - %v ", m)
		log.Print(err)
		return ""
	}
	log.Printf("Sending json - %s", debugDeidentify(buf))

	resp, err := http.Post(base, "application/json", bytes.NewReader(buf))
	if err != nil {
		log.Printf("Post failed to - %s: ", base)
		log.Print(err)
		return ""
	}
	defer resp.Body.Close()

	// TODO limit reader and timeout
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Response body failed - ")
		log.Print(err)
		return ""
	}
	log.Printf("Content rcv - %s", data)

	var tok TokenResp
	err = json.Unmarshal(data, &tok)
	if err != nil {
		log.Printf("Unmarshal failed on buffer - %s ", data)
		log.Print(err)
		return ""
	}

	return tok.AccessToken
}

func debugDeidentify(b []byte) string {
	re := regexp.MustCompile(`client_secret:[^,]*`)
	repl := []byte("client_secret: ***")

	result := re.ReplaceAll(b, repl)

	return string(result)
}
