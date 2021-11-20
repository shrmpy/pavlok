package ebs

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type HandlerFunc func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

/*
func MiddlewareTemplate(next HandlerFunc) HandlerFunc {
	return func(ev events.APIGatewayProxyRequest)
			(events.APIGatewayProxyResponse, error) {
		return next(ev)
	}
}
*/

func MiddlewareCORS(conf *Config, next HandlerFunc) HandlerFunc {
	return func(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// preflight check is short-circuited
		if ev.HTTPMethod == "OPTIONS" {
			return blankResponse(conf, "", http.StatusOK), nil
		}
		// without next, just act same as preflight
		if next == nil {
			return blankResponse(conf, "", http.StatusOK), nil
		}

		// run next handler along chain
		resp, err := next(ev)
		if err != nil {
			return resp, err
		}

		// post-process
		resp.Headers = enableCors(conf, resp.Headers)

		return resp, nil
	}
}

func enableCors(conf *Config, headers map[string]string) map[string]string {
	m := map[string]string{
		"Access-Control-Allow-Origin":  conf.Hostname(),
		"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT",
		"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	}
	// TODO merge, if CORS headers exist
	for key, val := range headers {
		m[key] = val
	}
	return m
}

func blankResponse(conf *Config, descr string, status int) events.APIGatewayProxyResponse {

	h := enableCors(conf, make(map[string]string))
	h["Content-Type"] = "application/json"

	return events.APIGatewayProxyResponse{

		Headers:    h,
		StatusCode: status,
	}
}
