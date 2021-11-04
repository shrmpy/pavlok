package main

import (

	//"errors"
	"log"
	"net/http"
	"regexp"
	"testing"

	//"github.com/aws/aws-lambda-go/lambda/handlertrace"
	"github.com/aws/aws-lambda-go/events"
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
	vanilla := newResponse("", http.StatusOK)
	// prep test data
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "OPTIONS",
		Body:       "",
	}
	// run the logic
	result, err := handler(req)
	// verify results
	assert.IsType(t, nil, err)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, vanilla.Headers, result.Headers)
}

func TestMissingAuthorization(t *testing.T) {
	// prep test data
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Body:       "",
	}
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
	wrongHeader := make(map[string]string)
	wrongHeader["Authorization"] = "Bearer WRONG-KEY"
	req := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Body:       "",
		Headers:    wrongHeader,
	}
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

func ignoreTimestamp(body string) string {
	// workaround to ignore unique key in the content
	re := regexp.MustCompile(`"id":"\d+"`)
	filtered := re.ReplaceAllString(body, "")
	return filtered
}
