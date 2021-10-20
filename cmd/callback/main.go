package main

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var helper *service

func init() {
	secret := os.Getenv("EXTENSION_SECRET")
	helper = newService(decodeSecret(secret))
}

func main() {
	lambda.Start(handler)
}

func handler(ev events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r, err := checkState(ev)
	if err != nil {
		return r, nil
	}

	addBroadcaster(ev)
	payload := formatJs()

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		Body:            payload,
		IsBase64Encoded: false,
	}, nil
}

// JS that closes the pop-up because the
// auth provider (pavlok) sent the browser
// here after completing auth via the redirect_uri
func formatJs() string {
	var js strings.Builder

	js.WriteString("<html>")
	js.WriteString("<script>")

	js.WriteString(" window.close();")

	js.WriteString("</script>")
	js.WriteString("</html>")

	return js.String()
}
