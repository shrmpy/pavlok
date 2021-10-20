package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var helper *service

func init() {
	secret := os.Getenv("EXTENSION_SECRET")
	helper = newService(decodeSecret(secret))
}

//  *API* flow
//      extract channel id from request auth header
//      then use it as reference key in fauna query
//      to retrieve the access token.
//      Call pavlok API with access token.

func main() {
	lambda.Start(handler)
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

	hdr := enableCors()
	hdr["Content-Type"] = "application/json"
	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            string(payload),
		Headers:         hdr,
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
