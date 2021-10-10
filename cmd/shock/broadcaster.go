package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
)

// extract channel ID from claims
func broadcaster(ev events.APIGatewayProxyRequest) string {

	accessor := core.RequestAccessor{}
	r, err := accessor.EventToRequest(ev)
	if err != nil {
		log.Print("Convert request failed")
		log.Print(err)
		return ""
	}

	token, err := helper.token(r)
	if err != nil {
		log.Print("Token extract failed")
		log.Print(err)
		return ""
	}

	claims, err := helper.claims(token)
	if err != nil {
		log.Print("Get claims failed")
		log.Print(err)
		return ""
	}

        if claims.ChannelID == "" {
                log.Print("Missing channel ID in claims")
		return ""
        }

	//todo enforce role for now during initial testing (until cheers wired-up)
	if claims.Role != "broadcaster" {
                log.Print("Role must be broadcaster during testing period")
		return ""
	}

        return claims.ChannelID
}

