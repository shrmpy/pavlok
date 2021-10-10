package main

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
)


// for the callback route to have the broadcaster channel ID
// we extract it from the claims structure and temporarily
// store it in the fauna db table (and use a reference ID
// that we know will carry across the redirect -- state field)
func reserveBroadcaster(ev events.APIGatewayProxyRequest) error {
        accessor := core.RequestAccessor{}
        req, err := accessor.EventToRequest(ev)
        if err != nil {
                log.Print("Convert event failed, broadcaster not reserved")
                return errors.New("convert event")
        }

        token, err := helper.token(req)
        if err != nil {
                log.Print("JWT extract failed ")
                return errors.New("jwt extract")
        }

	claims := getClaims(req)
	if claims.ChannelID == "" {
		log.Print("Claims is missing the channel ID ")
		////return errors.New("claims role")
	}
	if claims.Role != "broadcaster" {
		log.Print("Claims is not from the broadcaster role")
		return errors.New("claims role")
	}

	// store channel ID (with token as the reference key)
	secret := os.Getenv("FAUNADB_SECRET")
	pairs := NewPairs(secret, "Netlify")
	pairs.Insert(claims.ChannelID, token)

        return nil
}
