package main

import (
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	ebs "github.com/shrmpy/pavlok"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	request events.APIGatewayProxyRequest
	expect  string
	err     error
}

//todo structured cases
func TestHandlers(t *testing.T) {
	tests := []testCase{}
	/*		defaultTestCase(),
	}*/
	for i, tc := range tests {
		response, err := handler(tc.request)
		assert.IsType(t, tc.err, err)
		assert.Equal(t, tc.expect, response.Body)
		log.Printf("Test %d: %s", i, response.Body)
	}
}

func TestPreflight(t *testing.T) {
	// prep test data
	conf := ebs.NewConfig()
	expectMethods := "POST, GET, OPTIONS, PUT"
	expectHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
	expectOrigin := ""
	req := newTestRequest("OPTIONS")

	// run the handler logic
	result, err := ebs.MiddlewareCORS(conf, handler)(req)
	assert.IsType(t, nil, err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	// check for expected CORS
	assert.Equal(t, expectMethods, result.Headers["Access-Control-Allow-Methods"])
	assert.Equal(t, expectHeaders, result.Headers["Access-Control-Allow-Headers"])
	assert.Equal(t, expectOrigin, result.Headers["Access-Control-Allow-Origin"])
}

func TestMissingAuthorization(t *testing.T) {
	// prep test data
	req := newTestRequest("GET")

	// run the logic
	result, err := handler(req)
	// verify results
	assert.IsType(t, nil, err)
	assert.Equal(t, http.StatusOK, result.StatusCode)

	// case specific expectation
	missing := newResponse("Missing authorization header", http.StatusUnauthorized)
	expected := ignoreTimestamp(missing.Body)
	actual := ignoreTimestamp(result.Body)
	// inspect content
	assert.Equal(t, expected, actual)
}

func TestWrongExtensionSecret(t *testing.T) {
	// prepare data
	req := newTestRequest("GET")
	req.Headers["Authorization"] = "Bearer WRONG-KEY"

	// run handler logic
	result, err := handler(req)
	assert.IsType(t, nil, err)
	assert.Equal(t, http.StatusOK, result.StatusCode)

	// case specific expectation
	missing := newResponse("Wrong authorization header", http.StatusUnauthorized)
	expected := ignoreTimestamp(missing.Body)
	actual := ignoreTimestamp(result.Body)
	// inspect content
	assert.Equal(t, expected, actual)
}

func TestAccessControlAllowOrigin(t *testing.T) {
	// prepare data
	conf := ebs.NewConfig()
	conf.ExtensionId("HOSTNAME-TEST")
	expectMethods := "POST, GET, OPTIONS, PUT"
	expectHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
	expectOrigin := "https://HOSTNAME-TEST.ext-twitch.tv"
	req := newTestRequest("GET")

	// run handler logic
	result, err := ebs.MiddlewareCORS(conf, handler)(req)
	assert.IsType(t, nil, err)
	assert.Equal(t, http.StatusOK, result.StatusCode)

	// check for expected CORS
	assert.Equal(t, expectMethods, result.Headers["Access-Control-Allow-Methods"])
	assert.Equal(t, expectHeaders, result.Headers["Access-Control-Allow-Headers"])
	assert.Equal(t, expectOrigin, result.Headers["Access-Control-Allow-Origin"])
}

func newTestRequest(meth string) events.APIGatewayProxyRequest {
	// generic CORS request
	h := make(map[string]string)
	h["Content-Type"] = "application/json"

	return events.APIGatewayProxyRequest{
		HTTPMethod: meth,
		Body:       "",
		Headers:    h,
	}
}

func ignoreTimestamp(body string) string {
	// workaround to ignore unique key in the content
	re := regexp.MustCompile(`"id":"\d+"`)
	filtered := re.ReplaceAllString(body, "")
	return filtered
}
