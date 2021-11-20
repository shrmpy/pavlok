package main

import (
	"log"
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

//  *API* flow
//      extract channel id from request auth header
//      then use it as reference key in fauna query
//      to retrieve the access token.
//      Call pavlok API with access token.

func main() {
	lambda.Start(
		ebs.MiddlewareCORS(conf,
			handler,
		),
	)
}

func handler(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := preprocess(ev)
	if err != nil {
		return r, nil
	}

	refid := broadcaster(ev)
	if refid == "" {
		return newResponse("channel required", 500), nil
	}

	token := query(refid)
	status := callAPI(token)
	payload := notify(status)
	webhook(status)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(payload),

		IsBase64Encoded: false,
	}, nil
}

func query(refid string) string {
	log.Printf("Query by broadcaster - %s", refid)

	secret := os.Getenv("FAUNADB_SECRET")
	pairs := NewPairs(secret, "Netlify")
	token := pairs.Query(refid)

	return token
}
