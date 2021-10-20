package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
)

// store channel ID as reference key for pavlok access token
func addBroadcaster(ev events.APIGatewayProxyRequest) error {

	statetok := tokenFromState(ev)
	claims, err := helper.claims(statetok)
	if err != nil {
		return err
	}
	/*
		if claims.ChannelID == "" {
			return errors.New("channel")
		}
		if claims.Role != "broadcaster" {
			return errors.New("caster role")
		}
	*/

	// obtain pavlok token from auth code
	pavloktok := tradeToken(ev)

	secret := os.Getenv("FAUNADB_SECRET")
	pairs := NewPairs(secret, "Netlify")
	pairs.Update(pavloktok, claims.ChannelID)

	return nil
}
