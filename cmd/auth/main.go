package main

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ebs "github.com/shrmpy/pavlok"
)

var helper *service
var conf *ebs.Config

func init() {
	conf = ebs.NewConfig()
	secret := os.Getenv("EXTENSION_SECRET")
	helper = newService(decodeSecret(secret))
}

func main() {
	lambda.Start(
		ebs.MiddlewareCORS(conf,
			handler,
		),
	)
}

// *auth* flow part 1
//      initiate Pavlok OAuth login (authenticate request);
//      this is achieved by constructing the pavlok URL with the expected params
//      then return this in JSON for the frontend to launch in a browser tab.
//      A not very obvious part is managing the state field for the redirect
//      to our callback.

func handler(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := preprocess(ev)
	if err != nil {
		return r, nil
	}

	state := generateState(ev)

	buf, _ := json.Marshal(redirectHeaders(state))

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(buf),
	}, nil
}

// construct headers to redirect;
// here as a convenience for the frontend to
// avoid having to construct hyperlink with pavlok fields
func redirectHeaders(state string) map[string]string {
	cid := os.Getenv("PAVLOK_CLIENT_ID")
	redir := os.Getenv("PAVLOK_REDIRECT_URI")

	base := "https://app.pavlok.com/oauth/authorize?"
	params := url.Values{}
	params.Add("client_id", cid)
	params.Add("redirect_uri", redir)
	params.Add("response_type", "code")
	params.Add("state", state)

	m := make(map[string]string)
	m["Location"] = base + params.Encode()
	return m
}
