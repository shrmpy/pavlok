package main

import (

	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// preprocess for the callback route
func checkState(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// verify state before carrying out requests

	token := tokenFromState(ev)
        cl, err := helper.claims(token)

	if err != nil {
		log.Print("Bad OAuth state value")
                return newResponse("Bad OAuth state value", http.StatusUnauthorized), err
	}

	if cl.ChannelID == "" {
		log.Print("Missing channel ID from claims ")
                return newResponse("Missing channel", http.StatusUnauthorized), err
	}
	if cl.Role != "broadcaster" {
		log.Printf("Missing broadcaster role:  %s", cl.Role)
                return newResponse("Missing caster role", http.StatusUnauthorized), err
	}

	log.Printf("DEBUG claims - chan: %s, role: %s", cl.ChannelID, cl.Role)
        return events.APIGatewayProxyResponse{}, nil
}

func tokenFromState(ev events.APIGatewayProxyRequest) string {
	st := ev.QueryStringParameters["state"]
	log.Printf("Verifying state - %s", st)
	return decryptState(st)
}

