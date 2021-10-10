/*

*/
package main

import (
	"errors"
	"log"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
)

func preprocess(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
        if ev.HTTPMethod == "OPTIONS" {
                // answer preflight checks
                return newResponse("", http.StatusOK), errors.New("preflight step")
        }

        accessor := core.RequestAccessor{}
        req, err := accessor.EventToRequest(ev)
        if err != nil {
		log.Print("Convert event failed response ")
                return newResponse(err.Error(), http.StatusInternalServerError), err
        }

	token, err := helper.token(req)
	if err != nil {
		log.Print("JWT extract failed ")
                return newResponse(token, http.StatusUnauthorized), errors.New("jwt extract")
	}
	/*
        descr := helper.verifyJWT(token)
        if descr != "" {
		log.Print("JWT validate failed response ")
                return newResponse(descr, http.StatusUnauthorized), errors.New("jwt step")
        }*/

	return events.APIGatewayProxyResponse{}, nil
}

// make state value that functions like a file hash
func generateState(ev events.APIGatewayProxyRequest) string {
        accessor := core.RequestAccessor{}
        req, err := accessor.EventToRequest(ev)
        if err != nil {
		log.Print("Convert event failed, state not created")
                return ""
        }

	token, err := helper.token(req)
	if err != nil {
		log.Print("JWT extract failed ")
                return ""
	}
	// wrap token for use as state value
	return encryptState(token)
}


// include CORS in response header
func newResponse(descr string, status int) events.APIGatewayProxyResponse {
	buf := encodeError(descr, status)
	hdr := enableCors()
	hdr["Content-Type"] = "application/json"

	return events.APIGatewayProxyResponse{
		Body: string(buf),
		Headers: hdr,
		StatusCode: http.StatusOK,
	}
}

// format as close to JSON:API standard
func encodeError(descr string, status int) []byte {
	arr := newArray(descr, status)
	anon := struct {
		Errors []map[string]string `json:"errors"`
	}{Errors: arr}

	b, err := json.Marshal(anon)
	if err != nil {
		log.Print("JSON marshal failed for error -")
		log.Print(err)
		return []byte{}
	}
	return b
}

func newArray(descr string, status int) []map[string]string {
	str := strconv.Itoa(status)
	id := formatTimestamp()

	// item is actually a map
	item := map[string]string {
		"detail": descr,
		"status": str,
		"code": "DEMO-"+str,
		"id": id,
	}

	// array of objects that compose the top-level "errors" field
	return []map[string]string {item}
}

func formatTimestamp() string {
	// use UTC timestamp as error ID to reference
	// the moment of occurance (until errors are captured)

	n := time.Now().UnixNano()
	return strconv.FormatInt(n, 10)
}

// TODO review before open access
func enableCors() map[string]string {
        return map[string]string {
                "Access-Control-Allow-Origin": "*",
                "Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
                "Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
        }
}

